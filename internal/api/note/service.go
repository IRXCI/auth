package note

import (
	"github.com/IRXCI/auth/internal/service"
	desc "github.com/IRXCI/auth/pkg/auth"
)

type Implementation struct {
	desc.UnimplementedUserAPIServer
	noteServ service.NoteService
}

func NewImplementation(noteServ service.NoteService) *Implementation {
	return &Implementation{
		noteServ: noteServ,
	}
}
