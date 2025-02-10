package service

import (
	"context"

	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	SendCoin(ctx context.Context, params model.SendCoinParams) error
}

type UserService struct {
	repo UserRepository

	logger *logger.ZapLogger
}

func NewUserService(db UserRepository, l *logger.ZapLogger) *UserService {
	return &UserService{
		repo:   db,
		logger: l,
	}
}

func (s *UserService) AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, params.Username)
	if err != nil {
		return nil, err // TODO: update err handle
	}

	if err := compareHashAndPassword(user.Password, params.Password); err != nil {
		s.logger.Error("ValidateUserCredentials() -> compareHashAndPassword() request | error",
			zap.Any("params", params),
			zap.Error(err),
		)
		// пароли не совпали
		// либо некорректная длина хэша
		// отдаем неверный логин или пароль
		return nil, ErrInvalidLoginOrPassword
	}

	return user, nil
}

// TODO: REWRITE
func (s *UserService) GetUserInfo(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err // TODO: update err handle
	}

	return user, nil
}

func (s *UserService) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	if err := s.repo.SendCoin(ctx, params); err != nil {
		return err // TODO: update err handle
	}

	return nil
}
