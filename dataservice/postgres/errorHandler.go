package postgres

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func errorDecode(err error) string {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return "UniqueViolation"
		default:
			return "db error"
		}

	}

	return "unknown error"
}
