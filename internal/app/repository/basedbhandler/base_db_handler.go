package basedbhandler

import (
	"context"
)

type DBHandler interface {
	Execute(ctx context.Context, statement string, args ...interface{}) error
	ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error
	Query(ctx context.Context, statement string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, statement string, args ...interface{}) (Row, error)
	Close()
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type Row interface {
	Scan(dest ...interface{}) error
}
