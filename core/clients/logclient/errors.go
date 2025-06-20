package logclient

import (
	"errors"
	"fmt"
)

var (
	ErrClient       = errors.New("auth client")
	ErrBadRequest   = fmt.Errorf("%w: bad request", ErrClient)
	ErrUnauthorized = fmt.Errorf("%w: unauthorized", ErrClient)
	ErrForbidden    = fmt.Errorf("%w: forbidden", ErrClient)
	ErrNotFound     = fmt.Errorf("%w: not found", ErrClient)
	ErrConflict     = fmt.Errorf("%w: conflict", ErrClient)
	ErrInternal     = fmt.Errorf("%w: internal", ErrClient)
	ErrProcess      = fmt.Errorf("%w: process", ErrClient)
)
