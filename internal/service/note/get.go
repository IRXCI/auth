package note

import (
	"context"

	"github.com/IRXCI/auth/internal/model"
)

func (s *serv) GetUser(ctx context.Context, id int64) (*model.Note, error) {
	getUser, err := s.noteRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return getUser, nil
}
