package handler

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/api/convertor"

	desc "github.com/IRXCI/auth/pkg/auth"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) UpdateUser(ctx context.Context,
	info *desc.UpdateUserRequest) (*emptypb.Empty, error) {

	_, err := i.auth.UpdateUser(ctx, convertor.DomainFromDescUpdate(info))
	if err != nil {
		log.Printf("api failed to update user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
