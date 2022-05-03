package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
	"github.com/n0byk/loyalty/dataservice"
	"go.uber.org/zap"
)

var App Service
var AppConfig appConfig

type appConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"  envDefault:":8080"`
	DSN           string `env:"DATABASE_DSN" envDefault:"postgres://loyalty:loyalty@localhost:5432/loyalty?sslmode=disable"`
}

type Service struct {
	Logger  *zap.Logger
	Storage dataservice.Repository
}

func InitConfig(logger *zap.Logger) appConfig {

	flag.StringVar(&AppConfig.ServerAddress, "a", "localhost:8080", "RUN_ADDRESS")
	flag.StringVar(&AppConfig.DSN, "d", "", "DATABASE_URI")

	if err := env.Parse(&AppConfig); err != nil {
		logger.Error("Unset vars", zap.Error(err))
		os.Exit(1)
	}
	return AppConfig
}
