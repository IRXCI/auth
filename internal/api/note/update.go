package note

import (
	"context"

	"github.com/IRXCI/auth/internal/convertor"

	desc "github.com/IRXCI/auth/pkg/auth"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) UpdateUser(ctx context.Context,
	info *desc.UpdateUserRequest) (*emptypb.Empty, error) {

	_, err := i.noteServ.UpdateUser(ctx, convertor.ModelFromDescUpd(info))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
