package model

type User struct {
	ID       uint   `db:"id"`
	Username string `db:"username"`
	Password string `db:"password_hash"`
}
