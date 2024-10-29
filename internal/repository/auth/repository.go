package auth

import (
	"context"
	"errors"
	"log"
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

func validRole(info *domain.User) error {
	switch info.Role {
	case "ADMIN":
		return nil
	case "USER":
		return nil
	default:
		return errors.New("picked wrong role")
	}
}

func (r *repo) CreateUser(ctx context.Context,
	info *domain.User) (int64, error) {

	err := validRole(info)
	if err != nil {
		log.Printf("picked wrong role")
		return 0, err
	}

	builderCreateUser := sq.Insert(modelRepo.TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(modelRepo.NameColumn, modelRepo.EmailColumn, modelRepo.RoleColumn).
		Values(info.Name, info.Email, info.Role).
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

	return converter.ToAuthFromRepo(&auth), nil
}

func (r *repo) UpdateUser(ctx context.Context,
	info *domain.UserPlusId) (*emptypb.Empty, error) {

	builderUpdateUser := sq.Update(modelRepo.TableName).
		PlaceholderFormat(sq.Dollar).
		Set(modelRepo.UpdatedAtColumn, time.Now()).
		Where(sq.Eq{"id": info.Id})

	if info.Name != "" {
		builderUpdateUser = builderUpdateUser.Set(modelRepo.NameColumn, info.Name)
	}
	if info.Email != "" {
		builderUpdateUser = builderUpdateUser.Set(modelRepo.EmailColumn, info.Email)
	}

	if info.Role == "UNSPECIFIED" {
		log.Printf("picked wrong role")
		return nil, errors.New("picked wrong role")
	}

	if info.Role != "" {
		builderUpdateUser = builderUpdateUser.Set(modelRepo.RoleColumn, info.Role)
	}

	query, args, err := builderUpdateUser.ToSql()
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
