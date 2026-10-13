package group

import "errors"

var ErrInvalidID = errors.New("invalid id")
var ErrUserDontHavePermissionToView = errors.New("user don't have permission to view this item")
