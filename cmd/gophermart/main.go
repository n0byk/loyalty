package main

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/n0byk/loyalty/api/http/endpoints"
	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/dataservice/postgres"
	worker "github.com/n0byk/loyalty/workers"
)

func main() {
	logger := config.LoggerInit()
	appConfig := config.InitConfig(logger)

	pool, err := pgxpool.Connect(context.Background(), appConfig.DSN)
	if err != nil {
		logger.Error("Unable to connect to database: ", zap.Error(err))
		os.Exit(1)
	}
	storage := postgres.PostgresRepository(pool, logger)

	postgres.Migration(logger, appConfig.DSN)

	config.App = config.Service{Storage: storage, Logger: logger}

	go worker.AccrualAskWorker()

	logger.Info("ListenAndServe", zap.String("run_address", appConfig.ServerAddress))
	if err := http.ListenAndServe(appConfig.ServerAddress, endpoints.InitEndpoints(logger)); err != nil {
		logger.Error("ListenAndServe", zap.Error(err))
		os.Exit(1)
	}

}
