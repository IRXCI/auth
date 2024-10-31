package converter

import (
	"fmt"

	desc "github.com/IRXCI/auth/pkg/auth"
)

type Role string

const (
	Admin Role = "ADMIN"
	User  Role = "USER"
)

func RoleToDB(role string) (Role, error) {
	switch role {
	case desc.Role_ADMIN.String():
		return Admin, nil
	case desc.Role_USER.String():
		return User, nil
	default:
		return "", fmt.Errorf("failed to convert role %v to db", role)
	}
}
