package auth

import "errors"

var ErrInvalidUsername = errors.New("username must be provided")
var ErrInvalidPassword = errors.New("password must be at least 8 characters")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrUserNotFound = errors.New("user not found")
