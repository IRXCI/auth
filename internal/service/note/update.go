package note

import (
	"context"

	"github.com/IRXCI/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) UpdateUser(ctx context.Context, info *model.UserPlusId) (*emptypb.Empty, error) {
	_, err := s.noteRepo.UpdateUser(ctx, info)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
