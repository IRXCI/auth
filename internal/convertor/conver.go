package convertor

import (
	"github.com/IRXCI/auth/internal/model"
	desc "github.com/IRXCI/auth/pkg/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelFromDescCreate(info *desc.CreateUserRequest) *model.User {
	var res string
	switch info.UserAuth.Role {
	case desc.Role_USER:
		res = desc.Role_USER.String()
	case desc.Role_ADMIN:
		res = desc.Role_ADMIN.String()
	}

	return &model.User{
		Name:  info.UserAuth.Name,
		Email: info.UserAuth.Email,
		Role:  res,
	}
}

func DescGetResFromModel(note *model.Note) *desc.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        note.Id,
		UserAuth:  DescUserFromModel(note),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func DescUserFromModel(info *model.Note) *desc.User {
	var res desc.Role
	switch info.Role {
	case desc.Role_USER.String():
		res = desc.Role_USER
	case desc.Role_ADMIN.String():
		res = desc.Role_ADMIN
	}

	return &desc.User{
		Name:  info.Name,
		Email: info.Email,
		Role:  res,
	}
}

func ModelFromDescUpd(n *desc.UpdateUserRequest) *model.UserPlusId {
	return &model.UserPlusId{
		Id:    n.Id,
		Name:  n.Name.Value,
		Email: n.Email.Value,
		Role:  n.Role.String(),
	}
}
