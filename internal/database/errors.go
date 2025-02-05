package database

import (
	"errors"
)

var (
	ErrNoRecord            = errors.New("no records")
	ErrMultipleRecords     = errors.New("multiple records")
	ErrRecordAlreadyExists = errors.New("already exists")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrNegativeRowIndex    = errors.New("negative row index")
)
