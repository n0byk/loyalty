package postgres

import (
	"errors"

	"github.com/jackc/pgconn"
)

func errorDecode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return "010101"
}
