package database

import (
	"errors"
)

var ErrNoRecord = errors.New("no records")
var ErrMultipleRecords = errors.New("multiple records")
var ErrRecordAlreadyExists = errors.New("already exists")
var ErrNegativeRowIndex = errors.New("negative row index")
