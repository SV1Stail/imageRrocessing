package server

import (
	"errors"
)

var (
	ErrNoData         = errors.New("no data")
	ErrNoColor        = errors.New("no color")
	ErrWrongFormat    = errors.New("unsupported format")
	ErrNotImplemented = errors.New("not implemented")
)
