package note

import (
	"context"

	"github.com/IRXCI/auth/internal/convertor"
	desc "github.com/IRXCI/auth/pkg/auth"
)

func (i *Implementation) CreateUser(ctx context.Context,
	req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {

	id, err := i.noteServ.CreateUser(ctx, convertor.ModelFromDescCreate(req))
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
