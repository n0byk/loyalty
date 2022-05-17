package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/dataservice/models/entity"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

func (repository *postgreRepository) UpsertOrder(ctx context.Context, orderNumber string, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()
	_, err := repository.connection.Query(ctx, `insert into order_catalog (user_id, order_number, order_state) values ($1, $2, 'PROCESSING') on conflict do nothing returning user_id;`, userID, orderNumber)
	if err != nil {
		repository.logger.Error("UpsertOrder error", zap.Error(err))
		errorCode := errorDecode(err)
		return errorCode, errors.New(errorCode)
	}
	return userID, nil
}

func (repository *postgreRepository) SetOrder(ctx context.Context, orderNumber string, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()

	err := repository.connection.QueryRow(ctx, `insert into order_catalog (user_id, order_number, order_state) values ($1, $2, 'PROCESSING') returning user_id;`, userID, orderNumber).Scan(&userID)
	if err != nil {
		repository.logger.Error("SetOrder insert error", zap.Error(err))
		errorCode := errorDecode(err)
		return errorCode, errors.New(errorCode)
	}
	return userID, nil
}

func (repository *postgreRepository) CheckOrder(ctx context.Context, orderNumber string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()
	var userID string
	err := repository.connection.QueryRow(ctx, `select user_id from order_catalog where order_number = $1;`, orderNumber).Scan(&userID)

	if err != nil && err != pgx.ErrNoRows {
		repository.logger.Error("CheckOrder error", zap.Error(err))
		errorCode := errorDecode(err)
		return errorCode, errors.New(errorCode)
	}
	return userID, nil
}

func (repository *postgreRepository) GetOrder(ctx context.Context, userID string) ([]entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()

	var ordersSlice []entity.Order

	orders, err := repository.connection.Query(ctx, `select oc.order_number, oc.order_state, oc.update_time, coalesce(bc.accrue, 0) as accrual from order_catalog oc left join balance_catalog bc on bc.order_number = oc.order_number where user_id = $1;`, userID)
	if err != nil {
		errorCode := errorDecode(err)
		repository.logger.Error("GetOrder error", zap.Error(err))
		return ordersSlice, errors.New(errorCode)
	}

	for orders.Next() {
		var r entity.Order
		err := orders.Scan(&r.OrderNumber, &r.OrderState, &r.UpdateTime, &r.Accrual)
		if err != nil {
			repository.logger.Error("GetOrder error", zap.Error(err))

			return nil, err
		}
		ordersSlice = append(ordersSlice, r)
	}
	if err := orders.Err(); err != nil {
		repository.logger.Error("GetOrder error", zap.Error(err))

		return nil, err
	}
	return ordersSlice, nil
}

func (repository *postgreRepository) GetNewOrder(ctx context.Context) ([]entity.OrderIDNumber, string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()
	var orders []entity.OrderIDNumber

	rows, err := repository.connection.Query(ctx, `SELECT order_id, order_number FROM public.order_catalog WHERE order_state in('PROCESSING', 'NEW');`)

	if err != nil && err != pgx.ErrNoRows {
		repository.logger.Error("GetNewOrder error", zap.Error(err))
		errorCode := errorDecode(err)
		return orders, errorCode, errors.New(errorCode)
	}

	for rows.Next() {
		var r entity.OrderIDNumber
		err := rows.Scan(&r.OrderID, &r.OrderNumber)
		if err != nil {
			repository.logger.Error("GetWithdraws error", zap.Error(err))
			return orders, err.Error(), nil
		}
		orders = append(orders, r)
	}

	return orders, "", nil
}

func (repository *postgreRepository) SetOrderStatus(ctx context.Context, OrderID string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()
	_, err := repository.connection.Query(ctx, `update order_catalog set order_state = $1 where order_id = $2;`, status, OrderID)
	if err != nil {
		repository.logger.Error("SetOrderStatus error", zap.Error(err))
		errorCode := errorDecode(err)
		return errors.New(errorCode)
	}
	return nil
}

func (repository *postgreRepository) PostWithdraw(ctx context.Context, orderNumber string, sum float32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.DefaultCtxTimeout)
	defer cancel()

	err := repository.connection.QueryRow(ctx, `insert into balance_catalog (order_number, withdraw) values ($1, $2) returning order_number;`, orderNumber, sum).Scan(&orderNumber)
	if err != nil {
		repository.logger.Error("PostWithdraw error", zap.Error(err))
		errorCode := errorDecode(err)
		return errorCode, errors.New(errorCode)
	}
	return orderNumber, nil
}

func (repository *postgreRepository) PostAccrue(ctx context.Context, orderNumber string, sum float32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := repository.connection.QueryRow(ctx, `insert into balance_catalog (order_number, accrue) values ($1, $2) returning order_number;`, orderNumber, sum).Scan(&orderNumber)
	if err != nil {
		repository.logger.Error("PostAccrue error", zap.Error(err))
		errorCode := errorDecode(err)
		return errorCode, errors.New(errorCode)
	}
	return orderNumber, nil
}

func (repository *postgreRepository) GetWithdraws(ctx context.Context, userID string) ([]entity.OrderWithdrawals, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var ordersSlice []entity.OrderWithdrawals

	orders, err := repository.connection.Query(ctx, `select bc.order_number, bc.withdraw as withdraw, bc.update_time from balance_catalog bc left join order_catalog oc on oc.order_number = bc.order_number where oc.user_id = $1 and bc.withdraw is not null;`, userID)
	if err != nil {
		errorCode := errorDecode(err)
		repository.logger.Error("GetWithdraws error", zap.Error(err))
		return ordersSlice, errors.New(errorCode)
	}

	for orders.Next() {
		var r entity.OrderWithdrawals
		err := orders.Scan(&r.OrderID, &r.OrderSum, &r.ProcessedAt)
		if err != nil {
			repository.logger.Error("GetWithdraws error", zap.Error(err))
			return nil, err
		}
		ordersSlice = append(ordersSlice, r)
	}
	if err := orders.Err(); err != nil {
		repository.logger.Error("GetWithdraws error", zap.Error(err))
		return nil, err
	}
	return ordersSlice, nil
}
