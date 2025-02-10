package model

type AuthUserParams struct {
	Username string
	Password string
}

type SendCoinParams struct {
	FromUser uint
	ToUser   string
	Amount   int
}
