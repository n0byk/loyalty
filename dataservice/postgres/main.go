package postgres

import (
	"github.com/jackc/pgx/v4"
	"github.com/n0byk/loyalty/dataservice"
	"go.uber.org/zap"
)

type postgreRepository struct {
	connection *pgx.Conn
	logger     *zap.Logger
}

func PostgresRepository(db *pgx.Conn, logger *zap.Logger) dataservice.Repository {
	return &postgreRepository{
		connection: db,
		logger:     logger,
	}
}
