package database

import "errors"

var (
	ErrInvalidLoginOrPassword = errors.New("database: invalid login or password")
	ErrQueryFailed            = errors.New("database: query failed")
	ErrScanFailed             = errors.New("database: scan failed")
	ErrRowsFailed             = errors.New("database: rows failed")
	ErrNoFieldsToUpdate       = errors.New("database: no fields to update")

	ErrFailedToBeginTx         = errors.New("database: failed to begin transaction")
	ErrNotFound                = errors.New("database: not found")
	ErrFailedToFindRecipient   = errors.New("database: failed to find recipient")
	ErrFailedToFetchBalance    = errors.New("database: failed to fetch sender balance")
	ErrInsufficientFunds       = errors.New("database: insufficient funds")
	ErrFailedToDebitSender     = errors.New("database: failed to debit sender")
	ErrFailedToCreditRecipient = errors.New("database: failed to credit recipient")
	ErrFailedToSaveTx          = errors.New("database: failed to save transaction")
	ErrFailedToCommitTx        = errors.New("database: failed to commit transaction")
)
