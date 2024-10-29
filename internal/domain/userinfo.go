package domain

import (
	"database/sql"
	"time"
)

type UserInfo struct {
	Id        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
