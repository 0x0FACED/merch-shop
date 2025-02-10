package service

import (
	"errors"
)

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrQueryFailed            = errors.New("query failed")
	ErrScanFailed             = errors.New("scan failed")
	ErrRowsFailed             = errors.New("rows failed")
	ErrParseDate              = errors.New("error parse date")
)

var (
	ErrInvalidIDType = errors.New("invalid id type")
)

var (
	ErrNotFound = errors.New("not found")
)

var (
	ErrFailedToBeginTx         = errors.New("failed to begin tx")
	ErrInvalidRequest          = errors.New("invalid request")
	ErrFailedToFindRecipient   = errors.New("failed to find recipient")
	ErrFailedToFetchBalance    = errors.New("failed to fetch sender balance")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrFailedToDebitSender     = errors.New("failed to debit sender")
	ErrFailedToCreditRecipient = errors.New("failed to credit recipient")
	ErrFailedToSaveTransaction = errors.New("failed to save transaction")
	ErrFailedToCommitTx        = errors.New("failed to commit transaction")
)
