package service

import (
	"context"

	"github.com/0x0FACED/merch-shop/internal/model"
)

type userRepository interface {
	AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error)
	CreateUser(ctx context.Context, params model.CreateUserParams) (*model.User, error)
	GetUserInfo(ctx context.Context, params model.GetUserInfoParams) (*model.UserInfo, error)
	GetUserBalance(ctx context.Context, userID uint) (uint, error)
	SendCoin(ctx context.Context, params model.SendCoinParams) error
	BuyItem(ctx context.Context, params model.BuyItemParams) error
}
