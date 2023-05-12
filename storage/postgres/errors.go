package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/bohdanch-w/wheel/storage"
)

const (
	pgErrorCodeUniqueConstraint          = "23505"
	pgErrorForeignKeyConstraintViolation = "23503"
)

func ActualError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return storage.ErrRecordNotFound
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgErrorCodeUniqueConstraint:
			return storage.ErrUniqueConstraintViolation
		case pgErrorForeignKeyConstraintViolation:
			return storage.ErrForeignKeyConstraintViolation
		}
	}

	return err
}
