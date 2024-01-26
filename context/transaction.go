package context

import (
	"context"

	"github.com/google/uuid"
)

type TransactionContextKey struct{}

func WithTransactionID(c context.Context, t uuid.UUID) context.Context {
	return context.WithValue(c, TransactionContextKey{}, t)
}

func TransactionID(c context.Context) uuid.UUID {
	id, _ := c.Value(TransactionContextKey{}).(uuid.UUID)

	return id
}
