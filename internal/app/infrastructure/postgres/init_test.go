package postgres

import (
	ctx "context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"os"
	"testing"
)

const Datasource = "postgresql://practicum_ut:practicum_ut@127.0.0.1:5432/postgres"

var Log *zap.Logger
var target *PostgresqlHandlerTX

func initDatabase2() {
	conn, err := pgx.Connect(context.Background(), Datasource)
	if err != nil {
		fmt.Println("Can't init connection")
		panic(err)
	}
	_, err = conn.Exec(ctx.Background(), "create table if not exists test_table (a numeric, b varchar)")
	if err != nil {
		fmt.Println("Can't create test table")
		panic(err)
	}
	_, err = conn.Exec(ctx.Background(), "truncate table test_table")
	if err != nil {
		fmt.Println("Can't clear test table")
		panic(err)
	}
}

func TestMain(m *testing.M) {
	var err error
	Log, _ = zap.NewDevelopment()
	initDatabase2()
	target, err = NewPostgresqlHandlerTX(context.Background(), Datasource, Log)
	if err != nil {
		fmt.Println("can't init PostgresqlHandlerTX")
	}
	os.Exit(m.Run())
}
