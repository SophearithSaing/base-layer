package movie

import "errors"

var ErrInvalidID = errors.New("invalid id")
var ErrUserDontHavePermissionToView = errors.New("user don't have permission to view this item")
var ErrEmptyMovieBulk = errors.New("at least one movie must be provided")
var ErrInvalidRating = errors.New("rating must not be zero")
