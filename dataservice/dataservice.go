package dataservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0byk/loyalty/dataservice/models/entity"
)

type Repository interface {
	UserRegister(ctx context.Context, login, password, salt string) (uuid.UUID, string, error)
	UserLogin(ctx context.Context, login, password string) (entity.User, string, error)

	SetOrder(ctx context.Context, orderNumber string, userID string) (string, error)
	GetOrder(ctx context.Context, userID string) ([]entity.Order, error)

	GetBalance(ctx context.Context, userID string) (entity.Balance, error)

	CheckOrder(ctx context.Context, orderNumber string) (string, error)
	UpsertOrder(ctx context.Context, orderNumber string, userID string) (string, error)
	GetNewOrder(ctx context.Context) ([]entity.OrderIDNumber, string, error)
	SetOrderStatus(ctx context.Context, OrderID string, status string) error

	PostWithdraw(ctx context.Context, orderNumber string, sum float32) (string, error)
	GetWithdraws(ctx context.Context, userID string) ([]entity.OrderWithdrawals, error)

	PostAccrue(ctx context.Context, orderNumber string, sum float32) (string, error)

	DBPing() error
}
