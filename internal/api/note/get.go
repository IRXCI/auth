package note

import (
	"context"

	"github.com/IRXCI/auth/internal/convertor"

	desc "github.com/IRXCI/auth/pkg/auth"
)

func (i *Implementation) GetUser(ctx context.Context,
	id *desc.GetUserRequest) (*desc.GetUserResponse, error) {

	getUser, err := i.noteServ.GetUser(ctx, id.GetId())
	if err != nil {
		return nil, err
	}

	return convertor.DescGetResFromModel(getUser), nil
}
