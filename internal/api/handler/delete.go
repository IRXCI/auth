package handler

import (
	"context"
	"log"

	desc "github.com/IRXCI/auth/pkg/auth"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteUser(ctx context.Context,
	req *desc.DeleteUserRequest) (*emptypb.Empty, error) {

	_, err := i.auth.DeleteUser(ctx, req.Id)
	if err != nil {
		log.Printf("api failed to delete user: %v", err)
		return nil, err
	}

	log.Printf("handler: DeleteUser complete")

	return &emptypb.Empty{}, nil
}
