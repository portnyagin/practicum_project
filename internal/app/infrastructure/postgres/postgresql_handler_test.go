package postgres

import (
	ctx "context"
	"fmt"
	"golang.org/x/net/context"
	"testing"
)

const Datasource = "postgresql://practicum_ut:practicum_ut@127.0.0.1:5432/postgres"

func initDatabase() *PostgresqlHandler {
	h, err := NewPostgresqlHandler(ctx.Background(), Datasource)
	if err != nil {
		panic(err)
	}
	err = h.Execute(ctx.Background(), "create table if not exists test_table (a numeric, b varchar)")
	if err != nil {
		fmt.Println("Can't create test table")
		panic(err)
	}
	err = h.Execute(ctx.Background(), "truncate table test_table")
	if err != nil {
		fmt.Println("Can't clear test table")
		panic(err)
	}
	return h
}

func TestPostgresqlHandler_ExecuteBatch(t *testing.T) {
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
	handler := initDatabase()
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
			if err := handler.ExecuteBatch(ctx.Background(), tt.args.statement, paramArr); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresqlHandler_SelectRow(t *testing.T) {
	handler := initDatabase()
	row, err := handler.QueryRow(context.Background(), "select * from users where 1=0")
	var res int
	err = row.Scan(&res)
	fmt.Println(row)
	fmt.Println(err)
}
