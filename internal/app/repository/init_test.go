package repository

import (
	"context"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"os"
	"testing"
)

const Datasource = "postgresql://practicum_ut:practicum_ut@127.0.0.1:5432/postgres"

//var postgresHandler *infrastructure.PostgresqlHandler

func initDatabase(ctx context.Context, h *infrastructure.PostgresqlHandler) {
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
	//postgresHandler, err = infrastructure.NewPostgresqlHandler(context.Background(), Datasource)

	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
