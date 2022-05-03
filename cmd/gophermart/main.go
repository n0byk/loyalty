package main

import (
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/n0byk/loyalty/api"
	"github.com/n0byk/loyalty/api/http/endpoints"
	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/dataservice/postgres"
)

func main() {
	logger := api.LoggerInit()
	appConfig := config.InitConfig(logger)

	conn, err := pgx.Connect(context.Background(), appConfig.DSN)
	if err != nil {
		logger.Error("Unable to connect to database: ", zap.Error(err))
		os.Exit(1)
	}
	storage := postgres.PostgresRepository(conn, logger)
	defer conn.Close(context.Background())

	postgres.Migration(logger, appConfig.DSN)

	config.App = config.Service{Storage: storage, Logger: logger}

	logger.Info("ListenAndServe", zap.String("run_address", appConfig.ServerAddress))

	http.ListenAndServe(appConfig.ServerAddress, endpoints.InitEndpoints(logger))

}
