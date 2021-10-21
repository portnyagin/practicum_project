package app

import (
	"context"
	"fmt"
	cfg "github.com/portnyagin/practicum_project/internal/app/config"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/repository"
	"log"
)

func Start() {
	config := cfg.NewConfig()
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	postgresHandler, err := infrastructure.NewPostgresqlHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		fmt.Println("can't create postgres handler", err)
	}

	if config.Reinit {
		err = repository.ClearDatabase(context.Background(), postgresHandler)
		if err != nil {
			fmt.Println("can't clear database structure", err)
			return
		}
	}
	err = repository.InitDatabase(context.Background(), postgresHandler)
	if err != nil {
		fmt.Println("can't init database structure", err)
		return
	}

}
