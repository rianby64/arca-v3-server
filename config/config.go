package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	SpreadsheetID          string `env:"SPREADSHEET_ID,required"`
	ServiceCredentialsPath string `env:"SERVICE_CREDENTIALS_PATH,required"`
}

func showConfig(cfg *Config) {
	log.Printf("SPREADSHEET_ID=%s", cfg.SpreadsheetID)
	log.Printf("SERVICE_CREDENTIALS_PATH=%s", cfg.ServiceCredentialsPath)
}

func LoadConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	showConfig(cfg)

	return cfg
}
