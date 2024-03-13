package logger

import (
	"github.com/google/uuid"
)

func NewNilLogger() *NilLogger {
	return &NilLogger{}
}

type NilLogger struct{}

func (l *NilLogger) WithLevel(level LogLevel) Logger {
	return l
}

func (l *NilLogger) WithTransaction(id uuid.UUID) Logger {
	return l
}

func (l *NilLogger) WithError(err error) Logger {
	return l
}

func (l *NilLogger) With(key string, value any) Logger {
	return l
}

func (l *NilLogger) Debugf(msg string, args ...any) {
	// pass
}

func (l *NilLogger) Infof(msg string, args ...any) {
	// pass
}

func (l *NilLogger) Warnf(msg string, args ...any) {
	// pass
}

func (l *NilLogger) Errorf(msg string, args ...any) {
	// pass
}

func (l *NilLogger) Fatalf(msg string, args ...any) {
	// pass
}
