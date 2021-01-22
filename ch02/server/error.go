package server

import (
	"errors"
)

var (
	ErrEmpty      = errors.New("empty string")
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)
