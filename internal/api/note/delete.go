package note

import (
	"context"

	desc "github.com/IRXCI/auth/pkg/auth"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteUser(ctx context.Context,
	req *desc.DeleteUserRequest) (*emptypb.Empty, error) {

	_, err := i.noteServ.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
