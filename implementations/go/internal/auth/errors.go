package auth

import "errors"

var ErrUsernameNotProvided = errors.New("username must be provided")
var ErrPasswordTooShort = errors.New("password must be at least 8 characters")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrIncorrectPassword = errors.New("password is incorrect")

var ErrTokenNotFound = errors.New("token not found")
var ErrTokenIsRevoked = errors.New("token is revoked")
var ErrTokenIsExpired = errors.New("token is expired")
var ErrTokenIsInvalid = errors.New("token is invalid")
