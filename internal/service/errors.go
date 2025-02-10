package service

import "errors"

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrDatabaseInternal       = errors.New("internal database error")
	ErrParseDate              = errors.New("error parse date")
)

var (
	ErrInvalidIDType = errors.New("invalid id type")
)

var (
	ErrNotFound = errors.New("not found")
)
