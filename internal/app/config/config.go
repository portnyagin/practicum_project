package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"os"
)

type AppConfig struct {
	ServerAddress         string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseDSN           string `env:"DATABASE_URI" envDefault:"postgresql://practicum_project:practicum_project@127.0.0.1:5432/mdb"`
	AccrualServiceAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:3000"`
	Reinit                bool   `env:"REINIT" envDefault:"true"`
	ValidateOrderNum      bool   `env:"VALIDATE_ORDER" envDefault:"true"`
	EnableAccrual         bool   `env:"ENABLE_ACCRUAL" envDefault:"true"`
}

func (config *AppConfig) Init() error {
	fmt.Println(os.Args)
	if err := env.Parse(config); err != nil {
		fmt.Println("can't load service config", err)
		return err
	}

	pflag.StringVarP(&config.ServerAddress, "a", "a", config.ServerAddress, "Http-server address")
	pflag.StringVarP(&config.DatabaseDSN, "d", "d", config.DatabaseDSN, "Database connection string")
	pflag.StringVarP(&config.AccrualServiceAddress, "r", "r", config.AccrualServiceAddress, "Accrual Service Address")
	pflag.BoolVarP(&config.Reinit, "c", "c", config.Reinit, "Reinit database")
	pflag.BoolVarP(&config.ValidateOrderNum, "v", "v", config.ValidateOrderNum, "Validate order num")
	pflag.BoolVarP(&config.EnableAccrual, "y", "y", config.EnableAccrual, "Enable accrual processing")
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
