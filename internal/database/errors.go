package database

import "errors"

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrQueryFailed            = errors.New("query failed")
	ErrScanFailed             = errors.New("scan failed")
	ErrRowsFailed             = errors.New("rows failed")
	ErrNotFound               = errors.New("not found")
	ErrNoFieldsToUpdate       = errors.New("no fields to update")
	ErrTransactionFailed      = errors.New("transaction failed")
)
