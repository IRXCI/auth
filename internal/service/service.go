package service

import (
	"context"

	"github.com/IRXCI/auth/internal/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	CreateUser(ctx context.Context, info *domain.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*domain.UserInfo, error)
	UpdateUser(ctx context.Context, info *domain.UserWithId) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error)
}
