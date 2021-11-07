package basedbhandler

import (
	"context"
	"github.com/jackc/pgx/v4"
)

//go:generate mockgen -destination=mocks/mock_postgres_handler.go -package=mocks . DBHandler
type DBHandler interface {
	Execute(ctx context.Context, statement string, args ...interface{}) error
	ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error
	Query(ctx context.Context, statement string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, statement string, args ...interface{}) (Row, error)
	Close()
}

type TransactionalDBHandler interface {
	Execute(ctx context.Context, statement string, args ...interface{}) error
	ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error
	Query(ctx context.Context, statement string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, statement string, args ...interface{}) (Row, error)
	Commit(ctx context.Context)
	NewTx(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context)
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

//go:generate mockgen -destination=mocks/mock_row.go -package=mocks . Row
type Row interface {
	Scan(dest ...interface{}) error
}
