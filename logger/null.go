package logger

import "github.com/google/uuid"

type NullLogger struct{}

func (n *NullLogger) WithLevel(level LogLevel) Logger {
	return n
}

func (n *NullLogger) WithTransaction(id uuid.UUID) Logger {
	return n
}

func (n *NullLogger) WithError(err error) Logger {
	return n
}

func (n *NullLogger) With(key string, value any) Logger {
	return n
}

func (n *NullLogger) Debugf(msg string, args ...any) {
	// Do nothing
}

func (n *NullLogger) Infof(msg string, args ...any) {
	// Do nothing
}

func (n *NullLogger) Warnf(msg string, args ...any) {
	// Do nothing
}

func (n *NullLogger) Errorf(msg string, args ...any) {
	// Do nothing
}

func (n *NullLogger) Fatalf(msg string, args ...any) {
	// Do nothing
}
