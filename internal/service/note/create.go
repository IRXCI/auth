package note

import (
	"context"

	"github.com/IRXCI/auth/internal/model"
)

func (s *serv) CreateUser(ctx context.Context, info *model.User) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.noteRepo.CreateUser(ctx, info)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.noteRepo.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
