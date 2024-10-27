package converter

import (
	"github.com/IRXCI/auth/internal/model"
	modelRepo "github.com/IRXCI/auth/internal/repository/note/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {
	return &model.Note{
		Id:        note.Id,
		Name:      note.Name,
		Email:     note.Email,
		Role:      note.Role,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
