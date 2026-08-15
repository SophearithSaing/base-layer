package auth

import "errors"

var ErrUsernameNotProvided = errors.New("username must be provided")
var ErrPasswordTooShort = errors.New("password must be at least 8 characters")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrUserNotFound = errors.New("user not found")
