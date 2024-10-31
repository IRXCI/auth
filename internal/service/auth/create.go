package auth

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/domain"
)

func (a *AuthService) CreateUser(ctx context.Context, info *domain.User) (int64, error) {
	var id int64
	err := a.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = a.authRepo.CreateUser(ctx, info)
		if errTx != nil {
			return errTx
		}

		_, errTx = a.authRepo.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		log.Printf("service failed to create user: %v", err)
		return 0, err
	}

	return id, nil
}
