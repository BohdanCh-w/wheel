package context

import (
	"context"

	"github.com/bohdanch-w/wheel/logger"
)

type LoggerContextKey struct{}

func WithLogger(c context.Context, l logger.Logger) context.Context {
	return context.WithValue(c, LoggerContextKey{}, l)
}

func Logger(c context.Context) logger.Logger {
	l, ok := c.Value(LoggerContextKey{}).(logger.Logger)
	if !ok {
		return &logger.NullLogger{}
	}

	return l
}
