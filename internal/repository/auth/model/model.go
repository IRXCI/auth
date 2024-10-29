package model

import (
	"database/sql"
	"time"
)

const (
	TableName = "auth"

	IdColumn        = "id"
	NameColumn      = "name"
	EmailColumn     = "email"
	RoleColumn      = "role"
	CreatedAtColumn = "created_at"
	UpdatedAtColumn = "updated_at"
)

type UserInfo struct {
	Id        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Role      string       `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type User struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}

type UserPlusId struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}
