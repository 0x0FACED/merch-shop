package service

import (
	"context"
	"errors"

	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"go.uber.org/zap"
)

type userRepository interface {
	AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error)
	CreateUser(ctx context.Context, params model.CreateUserParams) (*model.User, error)
	SendCoin(ctx context.Context, params model.SendCoinParams) error
}

var _ userRepository = (*database.Postgres)(nil)

type UserService struct {
	repo userRepository

	logger *logger.ZapLogger
}

func NewUserService(db userRepository, l *logger.ZapLogger) *UserService {
	return &UserService{
		repo:   db,
		logger: l,
	}
}

func (s *UserService) AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error) {
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
				return nil, err
			}

			createParams := model.CreateUserParams{
				Username: params.Username,
				Password: hash,
			}

			user, err = s.repo.CreateUser(ctx, createParams)
			if err != nil {
				// Failed to create user, return
				return nil, err
			}
			// Successfully created user
			return user, nil
		}
		// ошибка базы
		return nil, err

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
		return nil, err
	}

	return user, nil
}

// TODO: REWRITE
func (s *UserService) GetUserInfo(ctx context.Context, username string) (*model.User, error) {
	//user, err := s.repo.AuthUserOrCreate(ctx, nil)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
	//return user, nil
}

func (s *UserService) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	if err := s.repo.SendCoin(ctx, params); err != nil {
		s.logger.Error("SendCoin() -> SendCoin() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		return err

	}

	return nil
}
