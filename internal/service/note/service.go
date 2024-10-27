package note

import (
	"github.com/IRXCI/auth/internal/client/db"
	"github.com/IRXCI/auth/internal/repository"
	"github.com/IRXCI/auth/internal/service"
)

type serv struct {
	noteRepo  repository.AuthRepository
	txManager db.TxManager
}

func NewService(noteRepo repository.AuthRepository, txManager db.TxManager) service.NoteService {
	return &serv{
		noteRepo:  noteRepo,
		txManager: txManager,
	}
}
