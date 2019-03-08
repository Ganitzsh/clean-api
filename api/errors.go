package api

import "errors"

var (
	ErrNotFound = errors.New("Record not found")
	ErrNilValue = errors.New("Cannot use nil value")
)
