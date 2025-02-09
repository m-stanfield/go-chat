package database

import (
	"errors"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrMultipleRecords     = errors.New("multiple records")
	ErrRecordAlreadyExists = errors.New("already exists")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrNegativeRowIndex    = errors.New("negative row index")
)
