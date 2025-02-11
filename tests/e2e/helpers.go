package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0x0FACED/merch-shop/internal/server"
	"github.com/stretchr/testify/assert"
)

// authUser отправляет запрос на регистрацию пользователя
func authUser(t *testing.T, username, password string, testServer *server.Server) string {
	reqBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	if testServer == nil {
		panic("nil test server")
	}
	testServer.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	json.Unmarshal(rec.Body.Bytes(), &response)
	return response["token"]
}
