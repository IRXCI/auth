package repository

import (
	"context"

	"github.com/IRXCI/auth/internal/repository/note/model"
	desc "github.com/IRXCI/auth/pkg/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, info *model.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*desc.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error)
}
