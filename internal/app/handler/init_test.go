package handler

import (
	"context"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure/postgres"
	"github.com/portnyagin/practicum_project/internal/app/repository"
	"github.com/portnyagin/practicum_project/internal/app/service"
	"go.uber.org/zap"
	"os"
	"testing"
)

var (
	log         *zap.Logger
	authService *service.AuthService
	auth        *Auth
)

const Datasource = "postgresql://practicum_ut:practicum_ut@127.0.0.1:5432/postgres"

func TestMain(m *testing.M) {
	//var err error
	log, _ = zap.NewDevelopment()
	postgresHandler, _ := postgres.NewPostgresqlHandler(context.Background(), Datasource)
	repository.ClearDatabase(context.Background(), postgresHandler)
	repository.InitDatabase(context.Background(), postgresHandler)
	repo, _ := repository.NewUserRepository(postgresHandler, log)
	authService = service.NewAuthService(repo, log)

	// TODO: в unit  тестах заменить моком
	auth = NewAuth("secret")

	//authHandler  = NewAuthHandler(, )
	os.Exit(m.Run())
}
