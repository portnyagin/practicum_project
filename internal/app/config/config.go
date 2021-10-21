package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"os"
)

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	DatabaseDSN   string `env:"DATABASE_DSN" envDefault:"postgresql://practicum:practicum@127.0.0.1:5432/postgres"`
	Reinit        bool   `env:"REINIT" envDefault:"true"`
}

func (config *AppConfig) Init() error {
	fmt.Println(os.Args)
	if err := env.Parse(config); err != nil {
		fmt.Println("can't load service config", err)
		return err
	}

	pflag.StringVarP(&config.ServerAddress, "a", "a", config.ServerAddress, "Http-server address")
	pflag.StringVarP(&config.DatabaseDSN, "d", "d", config.DatabaseDSN, "Database connection string")
	pflag.BoolVarP(&config.Reinit, "r", "r", config.Reinit, "Reinit database")
	pflag.Parse()

	if config.ServerAddress == "" || config.DatabaseDSN == "" {
		if err := env.Parse(&config); err != nil {
			fmt.Println("can't load service config", err)
			return err
		}
	}
	return nil
}

func NewConfig() *AppConfig {
	return &AppConfig{}
}
