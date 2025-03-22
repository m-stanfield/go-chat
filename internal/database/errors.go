package database

import (
	"errors"
)

// database errors
var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrMultipleRecords     = errors.New("multiple records")
	ErrRecordAlreadyExists = errors.New("already exists")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrNegativeRowIndex    = errors.New("negative row index")
)

// type conversion errors
var (
	ErrParsingValue             = errors.New("unable to parse value")
	ErrUnsupportedNegativeValue = errors.New("unsupported negative value")
)
