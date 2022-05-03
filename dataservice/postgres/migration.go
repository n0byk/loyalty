package postgres

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func Migration(logger *zap.Logger, dsn string) {
	m, err := migrate.New("file://dataservice/postgres/migrations", dsn)
	if err != nil {
		logger.Error("Migrations", zap.Error(err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		logger.Error("Migrations up error", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Migrations - done")

}
