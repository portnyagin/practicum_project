package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	cfg "github.com/portnyagin/practicum_project/internal/app/config"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure/postgres"
	"github.com/portnyagin/practicum_project/internal/app/repository"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/", mapRoute)

	err = http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		fmt.Println("can't start service")
		fmt.Println(err)
	}
}
