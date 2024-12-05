package helpers

import (
	"backend/internal/datasources/repositories"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

func PostgresErrorTransform(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pq.Error
	ok := errors.As(err, &pgErr)
	if ok {
		switch pgErr.Code {
		case "23505":
			return repositories.ErrorRowExists
		case "23503":
			return repositories.ErrorForeignKeyViolation
		}
	}

	// no rows in result set
	if errors.Is(err, sql.ErrNoRows) {
		return repositories.ErrorRowNotFound
	}

	return err
}
