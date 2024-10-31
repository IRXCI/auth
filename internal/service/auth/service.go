package auth

import (
	"github.com/IRXCI/auth/internal/client/db"
	"github.com/IRXCI/auth/internal/repository"
)

type AuthService struct {
	authRepo  repository.AuthRepository
	txManager db.TxManager
}

func NewService(authRepo repository.AuthRepository, txManager db.TxManager) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		txManager: txManager,
	}
}
