package model

type AuthUserParams struct {
	Username string
	Password string
}

type CreateUserParams struct {
	Username string
	Password string
}

type SendCoinParams struct {
	FromUser uint
	ToUser   string
	Amount   int
}

type GetUserInfoParams struct {
	ID uint
}

type BuyItemParams struct {
	UserID  uint
	Item    string
	Balance uint
}
