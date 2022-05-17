package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/dataservice/models/entity"
	"go.uber.org/zap"
)

func (repository *postgreRepository) GetBalance(ctx context.Context, userID string) (entity.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()
	var balance entity.Balance
	err := repository.connection.QueryRow(ctx, `select coalesce(SUM(bc.accrue), 0) - coalesce(SUM(bc.withdraw), 0) as accrue, coalesce(SUM(bc.withdraw), 0) as withdraw from balance_catalog bc left join order_catalog oc on oc.order_number = bc.order_number where oc.user_id = $1;`, userID).Scan(&balance.Current, &balance.Withdrawn)

	if err != nil && err != pgx.ErrNoRows {
		repository.logger.Error("GetBalance error", zap.Error(err))
		errorCode := errorDecode(err)
		return balance, errors.New(errorCode)
	}
	return balance, nil
}
