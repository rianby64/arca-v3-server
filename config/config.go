package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	version = "0.0.3 "
)

type Config struct {
	SpreadsheetID          string `env:"SPREADSHEET_ID,required"`
	ServiceCredentialsPath string `env:"SERVICE_CREDENTIALS_PATH,required"`

	ServerAddress string `env:"SERVER_ADDRESS,required"`
}

func showConfig(cfg *Config) {
	log.Printf("Configuration loaded: v%s", version)
	log.Printf("SPREADSHEET_ID\t\t= %s", cfg.SpreadsheetID)
	log.Printf("SERVICE_CREDENTIALS_PATH\t= %s", cfg.ServiceCredentialsPath)
	log.Printf("SERVER_ADDRESS\t\t= %s", cfg.ServerAddress)
}

func LoadConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	showConfig(cfg)

	return cfg
}
