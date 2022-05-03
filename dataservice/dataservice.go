package dataservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0byk/loyalty/dataservice/models/entity"
)

type Repository interface {
	UserRegister(ctx context.Context, login, password, salt string) (uuid.UUID, error)
	UserLogin(ctx context.Context, login, password string) (entity.User, error)
	// GetURL(ctx context.Context, key string) (string, error)
	// SetUserData(ctx context.Context, key, url, user string) error
	// GetUserData(ctx context.Context, user string) (string, error)
	DBPing() error
	// BulkDelete(ctx context.Context, urls []string, userID string) error
}
