package repository

import (
	"context"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure/postgres"
	"go.uber.org/zap"
	"os"
	"testing"
)

const Datasource = "postgresql://practicum_ut:practicum_ut@127.0.0.1:5432/postgres"

var postgresHandler *postgres.PostgresqlHandlerTX
var Log *zap.Logger

func initDatabase(ctx context.Context, h *postgres.PostgresqlHandlerTX) {
	err := ClearDatabase(ctx, h)
	if err != nil {
		fmt.Println("can't clear database")
		panic(err)
	}
	err = InitDatabase(ctx, h)
	if err != nil {
		fmt.Println("can't init database")
		panic(err)
	}
}

func TestMain(m *testing.M) {
	var err error

	Log, _ = zap.NewDevelopment()
	postgresHandler, err = postgres.NewPostgresqlHandlerTX(context.Background(), Datasource, Log)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
