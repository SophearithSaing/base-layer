package project

import "errors"

var ErrInvalidID = errors.New("invalid ID")
var ErrUserDontHavePermissionToView = errors.New("user don't have permission to view this item")
