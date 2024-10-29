package auth

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/domain"
)

func (a *AuthService) GetUser(ctx context.Context, id int64) (*domain.UserInfo, error) {
	getUser, err := a.authRepo.GetUser(ctx, id)
	if err != nil {
		log.Printf("service failed to get user: %v", err)
		return nil, err
	}

	return getUser, nil
}
