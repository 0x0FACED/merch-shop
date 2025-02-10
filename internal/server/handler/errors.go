package handler

import (
	"errors"
	"net/http"

	"github.com/0x0FACED/merch-shop/internal/service"
)

// MapServiceErrorToStatusCode маппит ошибки, полученные из Базы данных -> сервиса в статус коды
func MapServiceErrorToStatusCode(err error) int {
	switch {
	// 401 — Ошибки аутентификации
	case errors.Is(err, service.ErrInvalidLoginOrPassword):
		return http.StatusUnauthorized

	// 400 — Ошибки, связанные с неверными входными данными
	case errors.Is(err, service.ErrInsufficientFunds),
		errors.Is(err, service.ErrFailedToFindRecipient),
		errors.Is(err, service.ErrNotFound):
		return http.StatusBadRequest

	// 500 — Внутренние ошибки базы и транзакций
	case errors.Is(err, service.ErrQueryFailed),
		errors.Is(err, service.ErrScanFailed),
		errors.Is(err, service.ErrRowsFailed),
		errors.Is(err, service.ErrFailedToBeginTx),
		errors.Is(err, service.ErrFailedToFetchBalance),
		errors.Is(err, service.ErrFailedToDebitSender),
		errors.Is(err, service.ErrFailedToCreditRecipient),
		errors.Is(err, service.ErrFailedToSaveTransaction),
		errors.Is(err, service.ErrFailedToCommitTx):
		return http.StatusInternalServerError

	// 500 по дефолту
	default:
		return http.StatusInternalServerError
	}
}
