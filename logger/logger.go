package logger

import "github.com/google/uuid"

const TransactionKey = "transaction-id"

type Logger interface {
	WithLevel(level LogLevel) Logger
	WithTransaction(id uuid.UUID) Logger
	WithError(err error) Logger
	With(key string, value any) Logger
	Debugf(msg string, args ...any)
	Infof(msg string, args ...any)
	Warnf(msg string, args ...any)
	Errorf(msg string, args ...any)
	Fatalf(msg string, args ...any)
}
