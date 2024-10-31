package handler

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/api/convertor"

	desc "github.com/IRXCI/auth/pkg/auth"
)

func (i *Implementation) GetUser(ctx context.Context,
	id *desc.GetUserRequest) (*desc.GetUserResponse, error) {

	getUser, err := i.auth.GetUser(ctx, id.GetId())
	if err != nil {
		log.Printf("api failed to get user: %v", err)
		return nil, err
	}

	log.Printf("handler: GetUser complete")

	return convertor.DescGetResponseFromDomain(getUser), nil
}
