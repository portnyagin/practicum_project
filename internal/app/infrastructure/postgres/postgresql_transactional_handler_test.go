package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgresqlHandlerTX_getTx(t *testing.T) {
	baseContext := context.Background()
	tx, err := target.NewTx(baseContext)
	if err != nil {
		t.Error("can't init transaction")
	}
	ctxWithTransaction := context.WithValue(baseContext, basedbhandler.TransactionKey("tx"), tx)

	// В значении Нужен любой объект, отличный от TX
	ctxBad := context.WithValue(baseContext, basedbhandler.TransactionKey("tx"), errors.New("test"))
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "PostgresqlHandlerTX. getTx. Case #1. Positive",
			args: args{
				ctx: ctxWithTransaction,
			},
			wantErr: false,
		},
		{name: "PostgresqlHandlerTX. getTx. Case #2. Empty ",
			args: args{
				ctx: baseContext,
			},
			wantErr: true,
		},
		{name: "PostgresqlHandlerTX. getTx. Case #3. Bad transaction in context",
			args: args{
				ctx: ctxBad,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, err := target.getTx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPostgresqlHandlerTX_Transaction(t *testing.T) {
	type args struct {
		query  string
		commit bool
		rowCnt int
		resCnt int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "PostgresqlHandlerTX. Commit. Case #1.",
			args: args{
				query:  "insert into test_table (a,b) values ($1, $2)",
				commit: true,
				rowCnt: 3,
				resCnt: 3,
			},
			wantErr: false,
		},
		{name: "PostgresqlHandlerTX. Rollback. Case #2.",
			args: args{
				query:  "insert into test_table (a,b) values ($1, $2)",
				commit: false,
				rowCnt: 3,
				resCnt: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("init transaction")
			baseContext := context.Background()
			tx, err := target.NewTx(baseContext)
			if err != nil {
				t.Error("can't init transaction")
			}
			ctxWithTransaction := context.WithValue(baseContext, basedbhandler.TransactionKey("tx"), tx)

			fmt.Println("insert data undo transaction")
			for i := 0; i < tt.args.rowCnt; i++ {
				err = target.Execute(ctxWithTransaction, tt.args.query, i*10, tt.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if tt.args.commit {
				err = tx.Commit(ctxWithTransaction)
			} else {
				err = tx.Rollback(ctxWithTransaction)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Println("check result")
			rows, err := target.Query(context.Background(), "select * from test_table where b=$1", tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var cnt int
			for rows.Next() {
				cnt += 1
			}
			assert.Equal(t, tt.args.resCnt, cnt, "expectded %d rows, got %d", tt.args.resCnt, cnt)
		})
	}
}

func TestPostgresqlHandlerTX_ExecuteBatch(t *testing.T) {
	const batchStatement = "insert into test_table (a, b) values ($1, $2)"
	type args struct {
		statement string
		size      int
		a         []int
		b         []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ExecuteBatch test#1",
			args: args{statement: batchStatement,
				size: 0,
				a:    nil,
				b:    nil,
			},
			wantErr: false,
		},
		{name: "ExecuteBatch test#2",
			args: args{statement: batchStatement,
				size: 1,
				a:    []int{1},
				b:    []string{"str"},
			},
			wantErr: false,
		},
		{name: "ExecuteBatch test#3",
			args: args{statement: batchStatement,
				size: 3,
				a:    []int{1, 2, 3},
				b:    []string{"str1", "str2", "str3"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var paramArr [][]interface{}
			if tt.args.size > 0 {
				for i := 0; i < tt.args.size; i++ {
					var paramLine []interface{}
					paramLine = append(paramLine, tt.args.a[i])
					paramLine = append(paramLine, tt.args.b[i])
					paramArr = append(paramArr, paramLine)
				}
			}
			if err := target.ExecuteBatch(context.Background(), tt.args.statement, paramArr); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
