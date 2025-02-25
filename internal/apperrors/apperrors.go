package apperrors

import "errors"

var ErrConflict = errors.New("conflict: url already exists")
