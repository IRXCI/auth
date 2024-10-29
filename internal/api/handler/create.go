package handler

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/api/convertor"
	desc "github.com/IRXCI/auth/pkg/auth"
)

func (i *Implementation) CreateUser(ctx context.Context,
	req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {

	id, err := i.auth.CreateUser(ctx, convertor.DomainFromDescCreate(req))
	if err != nil {
		log.Printf("api failed to create user: %v", err)
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
