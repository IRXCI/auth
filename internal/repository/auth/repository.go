package auth

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/IRXCI/auth/internal/client/db"
	"github.com/IRXCI/auth/internal/domain"
	"github.com/IRXCI/auth/internal/repository"
	"github.com/IRXCI/auth/internal/repository/auth/converter"
	modelRepo "github.com/IRXCI/auth/internal/repository/auth/model"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context,
	info *domain.User) (int64, error) {

	dbRole, err := converter.RoleToDB(info.Role)
	if err != nil {
		slog.Error("failed to convert role to db", slog.Any("error", err))
		return 0, err
	}

	builderCreateUser := sq.Insert(modelRepo.TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(modelRepo.NameColumn, modelRepo.EmailColumn, modelRepo.RoleColumn).
		Values(info.Name, info.Email, dbRole).
		Suffix("RETURNING id")

	query, args, err := builderCreateUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "auth_repository.CreateUser",
		QueryRaw: query,
	}

	var UserID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&UserID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return 0, err
	}

	log.Printf("Inserted user with id: %d", UserID)

	return UserID, nil
}

func (r *repo) GetUser(ctx context.Context,
	id int64) (*domain.UserInfo, error) {

	builderGetUser := sq.Select(modelRepo.IdColumn, modelRepo.NameColumn,
		modelRepo.EmailColumn, modelRepo.RoleColumn, modelRepo.CreatedAtColumn, modelRepo.UpdatedAtColumn).
		From(modelRepo.TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.GetUser",
		QueryRaw: query,
	}

	var auth modelRepo.UserInfo
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(&auth.Id, &auth.Name, &auth.Email, &auth.Role, &auth.CreatedAt, &auth.UpdatedAt)
	if err != nil {
		log.Printf("failed to select user: %v", err)
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n",
		auth.Id, auth.Name, auth.Email, auth.Role, auth.CreatedAt, auth.UpdatedAt)

	return converter.ToDomainFromRepo(&auth), nil
}

func buildUpdatesMap(req *domain.UserWithId) (map[string]interface{}, bool) {
	updates := make(map[string]interface{})

	if email := req.Email; email != "" {
		updates[modelRepo.EmailColumn] = email
	}

	if name := req.Name; name != "" {
		updates[modelRepo.NameColumn] = name
	}

	if role, err := converter.RoleToDB(req.Role); err == nil {
		updates[modelRepo.RoleColumn] = role
	}

	if len(updates) != 0 {
		updates[modelRepo.UpdatedAtColumn] = time.Now()
		return updates, false
	}

	return updates, true
}

func (r *repo) UpdateUser(ctx context.Context,
	info *domain.UserWithId) (*emptypb.Empty, error) {

	updatedMap, noUpdates := buildUpdatesMap(info)
	if noUpdates {
		return &emptypb.Empty{}, nil
	}

	updateBuilder := sq.Update(modelRepo.TableName).
		SetMap(updatedMap).
		Where(sq.Eq{"id": info.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.UpdateUser",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

func (r *repo) DeleteUser(ctx context.Context,
	id int64) (*emptypb.Empty, error) {

	builderDeleteUser := sq.Delete(modelRepo.TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDeleteUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, err
	}

	log.Printf("User with id: %d, deleted", id)

	return &emptypb.Empty{}, nil
}
