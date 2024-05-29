package logger

import "github.com/google/uuid"

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

type NullLogger struct{}

func (n *NullLogger) WithLevel(_ LogLevel) Logger {
	return n
}

func (n *NullLogger) WithTransaction(_ uuid.UUID) Logger {
	return n
}

func (n *NullLogger) WithError(_ error) Logger {
	return n
}

func (n *NullLogger) With(_ string, _ any) Logger {
	return n
}

func (n *NullLogger) Debugf(_ string, _ ...any) {
	// Do nothing
}

func (n *NullLogger) Infof(_ string, _ ...any) {
	// Do nothing
}

func (n *NullLogger) Warnf(_ string, _ ...any) {
	// Do nothing
}

func (n *NullLogger) Errorf(_ string, _ ...any) {
	// Do nothing
}

func (n *NullLogger) Fatalf(_ string, _ ...any) {
	// Do nothing
}
