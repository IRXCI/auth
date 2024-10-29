package convertor

import (
	"github.com/IRXCI/auth/internal/domain"
	desc "github.com/IRXCI/auth/pkg/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func DomainFromDescCreate(info *desc.CreateUserRequest) *domain.User {
	var res string
	switch info.UserAuth.Role {
	case desc.Role_USER:
		res = desc.Role_USER.String()
	case desc.Role_ADMIN:
		res = desc.Role_ADMIN.String()
	}

	return &domain.User{
		Name:  info.UserAuth.Name,
		Email: info.UserAuth.Email,
		Role:  res,
	}
}

func DescGetResponseFromDomain(auth *domain.UserInfo) *desc.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if auth.UpdatedAt.Valid {
		updatedAt = timestamppb.New(auth.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        auth.Id,
		UserAuth:  DescUserFromDomain(auth),
		CreatedAt: timestamppb.New(auth.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func DescUserFromDomain(info *domain.UserInfo) *desc.User {
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

func DomainFromDescUpdate(n *desc.UpdateUserRequest) *domain.UserPlusId {
	return &domain.UserPlusId{
		Id:    n.Id,
		Name:  n.Name.Value,
		Email: n.Email.Value,
		Role:  n.Role.String(),
	}
}
