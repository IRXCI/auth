package model

type UserWithId struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}
