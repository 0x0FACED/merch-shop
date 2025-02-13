package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/internal/service/mocks"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testLogger() *logger.ZapLogger {
	cfg := config.LoggerConfig{
		LogLevel: "debug",
	}
	return logger.NewTestLogger(cfg)
}

// Тест успешной аутентификации
func TestAuthUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.AuthUserParams{Username: "testuser", Password: "test"}
	mockUser := &model.User{ID: 1, Username: "testuser", Password: "$2a$04$sjoS1Bf2A30VG0Vt0LkQf..KmKQCHuS5wvDG5RFTCQ2F1EtVVmTcm"}

	mockRepo.On("AuthUser", mock.Anything, params).Return(mockUser, nil)

	user, err := userService.AuthUser(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	mockRepo.AssertExpectations(t)
}

// Тест ошибки аутентификации
func TestAuthUser_Fail(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.AuthUserParams{Username: "unknown", Password: "password"}

	mockRepo.On("AuthUser", mock.Anything, params).Return(nil, errors.New("user not found"))

	user, err := userService.AuthUser(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthUser_CreateNewUser(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.AuthUserParams{Username: "newuser", Password: "password"}
	mockRepo.On("AuthUser", mock.Anything, params).Return(nil, database.ErrNotFound)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&model.User{ID: 2, Username: "newuser"}, nil)

	user, err := userService.AuthUser(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	mockRepo.AssertExpectations(t)
}

func TestAuthUser_CreateUserFail(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.AuthUserParams{Username: "testuser", Password: "test"}

	mockRepo.On("AuthUser", mock.Anything, params).Return(nil, database.ErrNotFound)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, database.ErrQueryFailed)

	_, err := userService.AuthUser(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrQueryFailed, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthUser_WrongHashFormat(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.AuthUserParams{Username: "testuser", Password: "test"}
	mockUser := &model.User{ID: 1, Username: "testuser", Password: "invalid-hash"}

	mockRepo.On("AuthUser", mock.Anything, params).Return(mockUser, nil)

	_, err := userService.AuthUser(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, service.ErrFailedComparingHashAndPassword, err)
	mockRepo.AssertExpectations(t)
}

// Тест успешного получения инфы о юзере
func TestGetUserInfo_Success(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.GetUserInfoParams{ID: 1}

	var expectedBalance uint = 1000

	mockRepo.On("GetUserInfo", mock.Anything, params).Return(&model.UserInfo{Coins: 1000}, nil)

	userInfo, err := userService.GetUserInfo(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, expectedBalance, userInfo.Coins)
	mockRepo.AssertExpectations(t)
}

// Тест ошибки получения инфы о юзере
func TestGetUserInfo_Fail(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.GetUserInfoParams{ID: 1}

	mockRepo.On("GetUserInfo", mock.Anything, params).Return(nil, database.ErrNotFound)

	userInfo, err := userService.GetUserInfo(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrNotFound, err)
	assert.Nil(t, userInfo)
	mockRepo.AssertExpectations(t)
}

// Тест успешной покупки предмета
func TestBuyItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.BuyItemParams{UserID: 1, Item: "hoody", Balance: 500}
	mockRepo.On("GetUserBalance", mock.Anything, params.UserID).Return(500, nil)
	mockRepo.On("BuyItem", mock.Anything, params).Return(nil)

	err := userService.BuyItem(context.Background(), params)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Тест, когда у пользователя не хватает денег на покупку
func TestBuyItem_NotEnoughBalance(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.BuyItemParams{UserID: 1, Item: "hoody", Balance: 100}
	mockRepo.On("GetUserBalance", mock.Anything, params.UserID).Return(100, nil)
	mockRepo.On("BuyItem", mock.Anything, params).Return(database.ErrInsufficientFunds)

	err := userService.BuyItem(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrInsufficientFunds, err)
	mockRepo.AssertExpectations(t)
}

// Тест, когда предмета не существует
func TestBuyItem_ItemDoesntExist(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.BuyItemParams{UserID: 1, Item: "not-exists-item", Balance: 100}
	mockRepo.On("GetUserBalance", mock.Anything, params.UserID).Return(100, nil)
	mockRepo.On("BuyItem", mock.Anything, params).Return(database.ErrNotFound)

	err := userService.BuyItem(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

// Тест покупки предмета, когда не получается получить баланс юзера
func TestBuyItem_GetUserBalanceError(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.BuyItemParams{UserID: 1, Item: "hoody"}
	mockRepo.On("GetUserBalance", mock.Anything, params.UserID).Return(0, database.ErrQueryFailed)

	err := userService.BuyItem(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrQueryFailed, err)
	mockRepo.AssertExpectations(t)
}

// Тест успешной передачи монеток
func TestSendCoin_Success(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "user2", Amount: 100}
	mockRepo.On("SendCoin", mock.Anything, params).Return(nil)

	err := userService.SendCoin(context.Background(), params)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Тест передачи монеток, если у пользователя недостаточно баланса
func TestSendCoin_NotEnoughBalance(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "user2", Amount: 5000}
	mockRepo.On("SendCoin", mock.Anything, params).Return(database.ErrInsufficientFunds)

	err := userService.SendCoin(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrInsufficientFunds, err)
	mockRepo.AssertExpectations(t)
}

// Тест передачи монеток несуществующему пользователю
func TestSendCoin_UserDoesntExist(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "not-exist-user", Amount: 50}
	mockRepo.On("SendCoin", mock.Anything, params).Return(database.ErrNotFound)

	err := userService.SendCoin(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, database.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestSendCoin_BeginTxFail(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "user2", Amount: 100}

	mockRepo.On("SendCoin", mock.Anything, params).Return(database.ErrFailedToBeginTx)

	err := userService.SendCoin(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, err, database.ErrFailedToBeginTx)
	mockRepo.AssertExpectations(t)
}

func TestSendCoin_FailedFetchBalance(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "user2", Amount: 100}

	mockRepo.On("SendCoin", mock.Anything, params).Return(database.ErrFailedToFetchBalance)

	err := userService.SendCoin(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, err, database.ErrFailedToFetchBalance)
	mockRepo.AssertExpectations(t)
}

func TestSendCoin_FailedCreditRecipient(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "user2", Amount: 100}

	mockRepo.On("SendCoin", mock.Anything, params).Return(database.ErrFailedToCreditRecipient)

	err := userService.SendCoin(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, err, database.ErrFailedToCreditRecipient)
	mockRepo.AssertExpectations(t)
}

func TestSendCoin_FailedDebitSender(t *testing.T) {
	mockRepo := new(mocks.MockMerchRepository)
	userService := service.NewUserService(mockRepo, testLogger())

	params := model.SendCoinParams{FromUser: 1, ToUser: "user2", Amount: 100}

	mockRepo.On("SendCoin", mock.Anything, params).Return(database.ErrFailedToDebitSender)

	err := userService.SendCoin(context.Background(), params)

	assert.Error(t, err)
	assert.Equal(t, err, database.ErrFailedToDebitSender)
	mockRepo.AssertExpectations(t)
}

func TestMapDBErrorToServiceError(t *testing.T) {
	err := service.MapDBErrorToServiceError(database.ErrInsufficientFunds)
	assert.Equal(t, service.ErrInsufficientFunds, err)

	err = service.MapDBErrorToServiceError(database.ErrNotFound)
	assert.Equal(t, service.ErrNotFound, err)

	err = service.MapDBErrorToServiceError(errors.New("unknown"))
	assert.Error(t, err)
}
