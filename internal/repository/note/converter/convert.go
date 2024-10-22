package converter

import (
	"github.com/IRXCI/auth/internal/repository/note/model"
	desc "github.com/IRXCI/auth/pkg/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromRepo(note *model.Note) *desc.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        note.Id,
		UserAuth:  ToUserFromRepo(note.UserNote),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserFromRepo(info *model.User) *desc.User {

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
