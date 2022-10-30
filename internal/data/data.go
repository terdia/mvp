package data

import (
	"errors"
)

var (
	ErrRecordNotFound       = errors.New("models: record not found")
	ErrDuplicateUsername    = errors.New("models: duplicate username")
	ErrInvalidCredentials   = errors.New("models: invalid credentials")
	ErrNoPermission         = errors.New("models: no permission")
	ErrDuplicateProductName = errors.New("models: you have created a product with the same name")
)

const (
	TokenScopeAuthentication = "authentication"
)
