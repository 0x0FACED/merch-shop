package service

import (
	"errors"
	"fmt"

	"github.com/0x0FACED/merch-shop/internal/database"
)

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrQueryFailed            = errors.New("query failed")
	ErrScanFailed             = errors.New("scan failed")
	ErrRowsFailed             = errors.New("rows failed")

	ErrFailedToBeginTx         = errors.New("failed to begin tx")
	ErrNotFound                = errors.New("not found")
	ErrFailedToFindRecipient   = errors.New("failed to find recipient")
	ErrFailedToFetchBalance    = errors.New("failed to fetch sender balance")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrFailedToDebitSender     = errors.New("failed to debit sender")
	ErrFailedToCreditRecipient = errors.New("failed to credit recipient")
	ErrFailedToSaveTransaction = errors.New("failed to save transaction")
	ErrFailedToCommitTx        = errors.New("failed to commit tx")

	ErrUnknown = errors.New("unknown error")
)

func MapDBErrorToServiceError(err error) error {
	switch {
	case errors.Is(err, database.ErrInvalidLoginOrPassword):
		return ErrInvalidLoginOrPassword

	case errors.Is(err, database.ErrInsufficientFunds):
		return ErrInsufficientFunds
	case errors.Is(err, database.ErrFailedToFindRecipient):
		return ErrFailedToFindRecipient
	case errors.Is(err, database.ErrNotFound):
		return ErrNotFound

	case errors.Is(err, database.ErrQueryFailed):
		return ErrQueryFailed
	case errors.Is(err, database.ErrScanFailed):
		return ErrScanFailed
	case errors.Is(err, database.ErrRowsFailed):
		return ErrRowsFailed
	case errors.Is(err, database.ErrFailedToBeginTx):
		return ErrFailedToBeginTx
	case errors.Is(err, database.ErrFailedToFetchBalance):
		return ErrFailedToFetchBalance
	case errors.Is(err, database.ErrFailedToDebitSender):
		return ErrFailedToDebitSender
	case errors.Is(err, database.ErrFailedToCreditRecipient):
		return ErrFailedToCreditRecipient
	case errors.Is(err, database.ErrFailedToSaveTransaction):
		return ErrFailedToSaveTransaction
	case errors.Is(err, database.ErrFailedToCommitTx):
		return ErrFailedToCommitTx

	default:
		return fmt.Errorf("%w: %w", ErrUnknown, err)
	}
}
