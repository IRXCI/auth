package auth

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *AuthService) UpdateUser(ctx context.Context, info *domain.UserPlusId) (*emptypb.Empty, error) {
	_, err := a.authRepo.UpdateUser(ctx, info)
	if err != nil {
		log.Printf("service failed to update user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
