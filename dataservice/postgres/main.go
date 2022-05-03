package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/n0byk/loyalty/dataservice"
	"go.uber.org/zap"
)

type postgreRepository struct {
	connection *pgxpool.Pool
	logger     *zap.Logger
}

func PostgresRepository(db *pgxpool.Pool, logger *zap.Logger) dataservice.Repository {
	return &postgreRepository{
		connection: db,
		logger:     logger,
	}
}
