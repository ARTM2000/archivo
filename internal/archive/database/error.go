package database

import "errors"

var (
	ErrUnhandled          = errors.New("unhandled error")
	ErrRecordNotFound     = errors.New("record not found")
	ErrDuplicateViolation = errors.New("duplication violence")
)
