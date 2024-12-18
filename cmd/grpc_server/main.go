package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	"github.com/IRXCI/auth/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/IRXCI/auth/pkg/auth"
	sq "github.com/Masterminds/squirrel"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	desc.UnimplementedUserAPIServer
	pool *pgxpool.Pool
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "../../.env", "path to config file")
}

func mapper(role string) (desc.Role, error) {
	var res desc.Role
	switch role {
	case desc.Role_USER.String():
		res = desc.Role_USER
	case desc.Role_ADMIN.String():
		res = desc.Role_ADMIN
	default:
		return res, errors.New("role cannot be converted")
	}
	return res, nil
}

func (s *server) CreateUser(ctx context.Context,
	req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {

	role := req.GetUserAuth().GetRole()
	if role == desc.Role_UNSPECIFIED {
		log.Printf("user picked wrong role")
		return nil, errors.New("picked wrong role")
	}

	builderCreateUser := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(req.GetUserAuth().GetName(),
			req.GetUserAuth().GetEmail(),
			req.GetUserAuth().GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderCreateUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var UserID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&UserID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return nil, err
	}

	log.Printf("Insert user with id: %d", UserID)

	return &desc.CreateUserResponse{
		Id: UserID,
	}, nil
}

func (s *server) GetUser(ctx context.Context,
	req *desc.GetUserRequest) (*desc.GetUserResponse, error) {

	builderGetUser := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderGetUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var (
		id                int64
		name, email, role string
		createdAt         time.Time
		updatedAt         sql.NullTime
	)

	err = s.pool.QueryRow(ctx, query, args...).
		Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select user: %v", err)
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n",
		id, name, email, role, createdAt, updatedAt)

	rolle, err := mapper(role)
	if err != nil {
		log.Printf("failed to converted role")
		return nil, err
	}

	return &desc.GetUserResponse{
		Id: req.GetId(),
		UserAuth: &desc.User{
			Name:  name,
			Email: email,
			Role:  rolle},

		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt.Time),
	}, nil
}

func (s *server) UpdateUser(
	ctx context.Context,
	req *desc.UpdateUserRequest,
) (*emptypb.Empty, error) {

	builderUpdateUser := sq.Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	if req.GetName().GetValue() != "" {
		builderUpdateUser = builderUpdateUser.Set("name", req.GetName().GetValue())
	}
	if req.GetEmail().GetValue() != "" {
		builderUpdateUser = builderUpdateUser.Set("email", req.GetEmail().GetValue())
	}
	if req.GetRole().String() != "" {
		builderUpdateUser = builderUpdateUser.Set("role", req.GetRole().String())
	}

	query, args, err := builderUpdateUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteUser(ctx context.Context,
	req *desc.DeleteUserRequest) (*emptypb.Empty, error) {

	builderDeleteUser := sq.Delete("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDeleteUser.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, err
	}

	log.Printf("User with id: %d, deleted", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIServer(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
