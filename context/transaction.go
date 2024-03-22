package context

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey uint8

const (
	TransactionIDKey ctxKey = iota
)

func WithTransactionID(c context.Context, t uuid.UUID) context.Context {
	return context.WithValue(c, TransactionIDKey, t)
}

func TransactionID(c context.Context) uuid.UUID {
	id, _ := c.Value(TransactionIDKey).(uuid.UUID)

	return id
}
