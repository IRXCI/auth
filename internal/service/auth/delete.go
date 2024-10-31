package auth

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *AuthService) DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error) {
	_, err := a.authRepo.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("service failed to delete user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
