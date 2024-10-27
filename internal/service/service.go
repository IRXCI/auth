package service

import (
	"context"

	"github.com/IRXCI/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NoteService interface {
	CreateUser(ctx context.Context, info *model.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.Note, error)
	UpdateUser(ctx context.Context, info *model.UserPlusId) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error)
}
