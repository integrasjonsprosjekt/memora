package models

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser = errors.New("invalid user data")
)
