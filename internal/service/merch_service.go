package service

import (
	"context"
	"errors"

	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"go.uber.org/zap"
)

var _ merchRepository = (*database.Postgres)(nil)

type MerchService struct {
	repo merchRepository

	logger *logger.ZapLogger
}

func NewUserService(db merchRepository, l *logger.ZapLogger) *MerchService {
	return &MerchService{
		repo:   db,
		logger: l,
	}
}

func (s *MerchService) AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error) {
	user, err := s.repo.AuthUser(ctx, params)
	if err != nil {
		s.logger.Error("AuthUser() -> AuthUser() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		// Не нашли юзера, значит создаем его
		if errors.Is(err, database.ErrNotFound) {
			// create user
			hash, err := calcHash(params.Password)
			if err != nil {
				return nil, MapDBErrorToServiceError(err)
			}

			createParams := model.CreateUserParams{
				Username: params.Username,
				Password: hash,
			}

			user, err = s.repo.CreateUser(ctx, createParams)
			if err != nil {
				// Failed to create user, return
				return nil, MapDBErrorToServiceError(err)
			}
			// Successfully created user
			return user, nil
		}
		// ошибка базы
		return nil, MapDBErrorToServiceError(err)
	}

	// юзер существует, проверяем пароль и хэш базы
	if err := compareHashAndPassword(user.Password, params.Password); err != nil {
		s.logger.Error("AuthUser() -> compareHashAndPassword() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		// пароли не совпали
		// либо некорректная длина хэша
		// отдаем неверный логин или пароль
		return nil, ErrFailedComparingHashAndPassword
	}

	return user, nil
}

// TODO: REWRITE
func (s *MerchService) GetUserInfo(ctx context.Context, params model.GetUserInfoParams) (*model.UserInfo, error) {
	userInfo, err := s.repo.GetUserInfo(ctx, params)
	if err != nil {
		s.logger.Error("GetUserInfo() -> GetUserInfo() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		return nil, MapDBErrorToServiceError(err)
	}

	return userInfo, nil
}

func (s *MerchService) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	if err := s.repo.SendCoin(ctx, params); err != nil {
		s.logger.Error("SendCoin() -> SendCoin() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		return MapDBErrorToServiceError(err)
	}

	return nil
}

func (s *MerchService) BuyItem(ctx context.Context, params model.BuyItemParams) error {
	balance, err := s.repo.GetUserBalance(ctx, params.UserID)
	if err != nil {
		s.logger.Error("BuyItem() -> GetUserBalance() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		return MapDBErrorToServiceError(err)
	}

	params.Balance = balance

	if err := s.repo.BuyItem(ctx, params); err != nil {
		s.logger.Error("BuyItem() -> BuyItem() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		return MapDBErrorToServiceError(err)
	}

	return nil
}
