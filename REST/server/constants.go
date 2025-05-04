package server

import (
	"errors"
)

var (
	ErrNotAllowed = errors.New("not allowed")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal error")
)
