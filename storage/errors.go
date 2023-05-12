package storage

import wherr "github.com/bohdanch-w/wheel/errors"

const (
	ErrRecordNotFound                = wherr.Error("record not found")
	ErrUniqueConstraintViolation     = wherr.Error("unique constraint violation")
	ErrForeignKeyConstraintViolation = wherr.Error("foreign key constraint violation")
)
