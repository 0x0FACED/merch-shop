package mocks

import (
	"context"

	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error) {
	args := m.Called(ctx, params)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) CreateUser(ctx context.Context, params model.CreateUserParams) (*model.User, error) {
	args := m.Called(ctx, params)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) GetUserInfo(ctx context.Context, params model.GetUserInfoParams) (*model.UserInfo, error) {
	args := m.Called(ctx, params)
	if userInfo, ok := args.Get(0).(*model.UserInfo); ok {
		return userInfo, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) GetUserBalance(ctx context.Context, userID uint) (uint, error) {
	args := m.Called(ctx, userID)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockMerchRepository) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockMerchRepository) BuyItem(ctx context.Context, params model.BuyItemParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}
