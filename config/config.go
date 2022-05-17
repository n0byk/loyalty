package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/n0byk/loyalty/dataservice"
	"go.uber.org/zap"
)

var App Service
var AppConfig appConfig

type appConfig struct {
	ServerAddress         string        `env:"RUN_ADDRESS"  envDefault:":8081"`
	DSN                   string        `env:"DATABASE_URI" envDefault:"postgres://loyalty:loyalty@localhost:5432/loyalty?sslmode=disable"`
	AccrualSystemAddress  string        `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
	DefaultCtxTimeout     time.Duration `env:"DEFAULT_CTX_TIMEOUT"`
	AccrualAskWorkerCount int           `env:"ACCRUAL_ASK_WORKER_COUNT"`
}

type Service struct {
	Logger  *zap.Logger
	Storage dataservice.Repository
}

func InitConfig(logger *zap.Logger) appConfig {

	flag.StringVar(&AppConfig.ServerAddress, "a", ":8081", "RUN_ADDRESS")
	flag.StringVar(&AppConfig.DSN, "d", "", "DATABASE_URI")
	flag.StringVar(&AppConfig.AccrualSystemAddress, "r", "", "ACCRUAL_SYSTEM_ADDRESS")
	flag.DurationVar(&AppConfig.DefaultCtxTimeout, "c", 2*time.Second, "DEFAULT_CTX_TIMEOUT")
	flag.IntVar(&AppConfig.AccrualAskWorkerCount, "w", 1, "ACCRUAL_ASK_WORKER_COUNT")

	if err := env.Parse(&AppConfig); err != nil {
		logger.Error("Unset vars", zap.Error(err))
		os.Exit(1)
	}
	flag.Parse()
	return AppConfig
}
