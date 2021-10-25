package app

import (
	"context"
	cfg "github.com/portnyagin/practicum_project/internal/app/config"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure/postgres"
	"github.com/portnyagin/practicum_project/internal/app/repository"
	"go.uber.org/zap"
	"log"
)

func Start() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	config := cfg.NewConfig()
	err = config.Init()
	if err != nil {
		logger.Fatal("can't init configuration", zap.Error(err))
	}

	postgresHandler, err := postgres.NewPostgresqlHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		logger.Fatal("can't create postgres handler", zap.Error(err))
	}

	if config.Reinit {
		err = repository.ClearDatabase(context.Background(), postgresHandler)
		if err != nil {
			logger.Fatal("can't clear database structure", zap.Error(err))
			return
		}
	}
	err = repository.InitDatabase(context.Background(), postgresHandler)
	if err != nil {
		logger.Fatal("can't init database structure", zap.Error(err))
		return
	}

}
