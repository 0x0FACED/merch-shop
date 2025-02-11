package handler

import "github.com/0x0FACED/merch-shop/internal/model"

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type InfoResponse struct {
	Coins       uint              `json:"coins"`
	Inventory   model.Inventory   `json:"inventory"`
	CoinHistory model.CoinHistory `json:"coinHistory"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
