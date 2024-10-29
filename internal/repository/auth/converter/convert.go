package converter

import (
	"github.com/IRXCI/auth/internal/domain"
	modelRepo "github.com/IRXCI/auth/internal/repository/auth/model"
)

func ToAuthFromRepo(auth *modelRepo.UserInfo) *domain.UserInfo {
	return &domain.UserInfo{
		Id:        auth.Id,
		Name:      auth.Name,
		Email:     auth.Email,
		Role:      auth.Role,
		CreatedAt: auth.CreatedAt,
		UpdatedAt: auth.UpdatedAt,
	}
}
