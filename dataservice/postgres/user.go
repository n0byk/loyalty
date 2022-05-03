package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/n0byk/loyalty/dataservice/models/entity"
	"go.uber.org/zap"
)

func (repository *postgreRepository) UserRegister(ctx context.Context, login, password, salt string) (uuid.UUID, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var userID uuid.UUID

	err := repository.connection.QueryRow(ctx, `insert into user_catalog (user_login, user_password, user_salt) values ($1, $2, $3) returning user_id;`, login, password, salt).Scan(&userID)
	if err != nil {
		repository.logger.Error("UserRegister insert error", zap.Error(err))
		errorCode := errorDecode(err)
		return userID, errorCode, errors.New(errorCode)
	}
	return userID, "", nil
}

func (repository *postgreRepository) UserLogin(ctx context.Context, login, password string) (entity.User, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var user entity.User

	err := repository.connection.QueryRow(ctx, `select uc.user_id, uc.user_login, uc.user_password, uc.user_salt from user_catalog uc where uc.user_login = $1 and uc.delete_time is null;`, login).Scan(&user.UserID, &user.UserLogin, &user.UserPassword, &user.UserSalt)
	if err != nil {
		repository.logger.Error("UserLogin insert error", zap.Error(err))
		errorCode := errorDecode(err)
		return user, errorCode, errors.New(errorCode)
	}
	return user, "", nil
}

func (repository *postgreRepository) DBPing() error {
	return repository.connection.Ping(context.Background())
}
