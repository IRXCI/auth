package model

type User struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}
