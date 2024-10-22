package note

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/IRXCI/auth/internal/repository"
	"github.com/IRXCI/auth/internal/repository/note/converter"
	"github.com/IRXCI/auth/internal/repository/note/model"
	desc "github.com/IRXCI/auth/pkg/auth"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	tableName = "auth"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context,
	info *model.User) (int64, error) {

	role := info.Role
	if role == "UNSPECIFIED" {
		log.Printf("picked wrong role")
		return 0, errors.New("picked wrong role")
	}

	builderCreateUser := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn).
		Values(info.Name, info.Email, info.Role).
		Suffix("RETURNING id")

	query, args, err := builderCreateUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, err
	}

	var UserID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&UserID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return 0, err
	}

	log.Printf("Insert user with id: %d", UserID)

	return UserID, nil
}

func (r *repo) GetUser(ctx context.Context,
	id int64) (*desc.GetUserResponse, error) {

	builderGetUser := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var note model.Note
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&note.Id, &note.UserNote, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		log.Printf("failed to select user: %v", err)
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n",
		note.Id, note.UserNote.Name, note.UserNote.Email, note.UserNote.Role, note.CreatedAt, note.UpdatedAt)

	return converter.ToNoteFromRepo(&note), nil
}

func (r *repo) UpdateUser(ctx context.Context,
	req *desc.UpdateUserRequest,
) (*emptypb.Empty, error) {

	builderUpdateUser := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	if req.GetName().GetValue() != "" {
		builderUpdateUser = builderUpdateUser.Set(nameColumn, req.GetName().GetValue())
	}
	if req.GetEmail().GetValue() != "" {
		builderUpdateUser = builderUpdateUser.Set(emailColumn, req.GetEmail().GetValue())
	}
	if req.GetRole().String() != "" {
		builderUpdateUser = builderUpdateUser.Set(roleColumn, req.GetRole().String())
	}

	query, args, err := builderUpdateUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

func (r *repo) DeleteUser(ctx context.Context,
	id int64) (*emptypb.Empty, error) {

	builderDeleteUser := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDeleteUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, err
	}

	log.Printf("User with id: %d, deleted", id)

	return &emptypb.Empty{}, nil
}
