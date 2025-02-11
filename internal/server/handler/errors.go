package handler

import (
	"errors"
	"net/http"

	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/service"
)

// MapServiceErrorToStatusCode маппит ошибки, полученные из базы в статус коды
func MapServiceErrorToStatusCode(err error) int {
	switch {
	// 401 — Ошибки аутентификации
	case errors.Is(err, database.ErrInvalidLoginOrPassword),
		errors.Is(err, service.ErrFailedComparingHashAndPassword):
		return http.StatusUnauthorized

	// 400 — Ошибки, связанные с неверными входными данными
	case errors.Is(err, database.ErrInsufficientFunds),
		errors.Is(err, database.ErrFailedToFindRecipient),
		errors.Is(err, database.ErrNotFound):
		return http.StatusBadRequest

	// 500 — Внутренние ошибки базы и транзакций
	case errors.Is(err, database.ErrQueryFailed),
		errors.Is(err, database.ErrScanFailed),
		errors.Is(err, database.ErrRowsFailed),
		errors.Is(err, database.ErrFailedToBeginTx),
		errors.Is(err, database.ErrFailedToFetchBalance),
		errors.Is(err, database.ErrFailedToDebitSender),
		errors.Is(err, database.ErrFailedToCreditRecipient),
		errors.Is(err, database.ErrFailedToSaveTransaction),
		errors.Is(err, database.ErrFailedToCommitTx):
		return http.StatusInternalServerError

	// 500 по дефолту
	default:
		return http.StatusInternalServerError
	}
}
