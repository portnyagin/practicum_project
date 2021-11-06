package model

import (
	"errors"
	"github.com/jackc/pgerrcode"
)

type DatabaseError struct {
	Err  error
	Code string
}

func (t *DatabaseError) Error() string {
	return t.Err.Error()
}

var (
	UniqueViolation DatabaseError = DatabaseError{Code: pgerrcode.UniqueViolation}
	NoRowFound      DatabaseError = DatabaseError{Err: errors.New("no rows in result set")}
)
