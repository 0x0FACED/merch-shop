package handler

type AuthRequest struct {
	Username string `json:"username" validate:"required,alphanum,max=255"`
	Password string `json:"password" validate:"required,alphanum,min=4,max=128"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required,alphanum,max=255"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}
