package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuyItem_Success проверяет покупку товара
func TestBuyItem_Success(t *testing.T) {
	token := authUser(t, "testuser", "password", testServer)

	req := httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "The purchase failed")
}

// TestBuyItem_NotEnoughMoney проверяет случай, когда не хватает монет
func TestBuyItem_NotEnoughMoney(t *testing.T) {
	token := authUser(t, "pooruser", "password", testServer)

	reqBody, _ := json.Marshal(map[string]any{
		"toUser": "testuser",
		"amount": 900,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	req = httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec = httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code, "There must be a waiver due to lack of funds")
}

// TestSendCoin_Success проверяет отправку монет
func TestSendCoin_Success(t *testing.T) {
	token := authUser(t, "sender", "password", testServer)

	reqBody, _ := json.Marshal(map[string]any{
		"toUser": "testuser",
		"amount": 100,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "The coin sending failed")
}

// TestSendCoin_NotEnoughMoney проверяет отказ при нехватке монет
func TestSendCoin_NotEnoughMoney(t *testing.T) {
	token := authUser(t, "poor_sender", "password", testServer)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"toUser": "testuser",
		"amount": 5000, // Больше, чем есть на балансе
	})

	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	fmt.Println(rec.Body)
	assert.Equal(t, http.StatusBadRequest, rec.Code, "There must be a waiver due to lack of funds")
}

// TestAuth_WrongPassword проверяет авторизацию с неверным паролем
func TestAuth_WrongPassword(t *testing.T) {
	authUser(t, "wrongpassuser", "password", testServer)

	reqBody, _ := json.Marshal(map[string]string{
		"username": "wrongpassuser",
		"password": "wrongpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "There should be a waiver of authorization")
}

// TestUserInfo_Success проверяет получение инфы о пользователе
func TestUserInfo_Success(t *testing.T) {
	token := authUser(t, "info_user", "password", testServer)

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	testServer.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Retrieving user information failed")
}
