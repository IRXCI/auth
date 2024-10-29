package handler

import (
	"github.com/IRXCI/auth/internal/service"
	desc "github.com/IRXCI/auth/pkg/auth"
)

type Implementation struct {
	desc.UnimplementedUserAPIServer
	auth service.AuthService
}

func NewImplementation(auth service.AuthService) *Implementation {
	return &Implementation{
		auth: auth,
	}
}
