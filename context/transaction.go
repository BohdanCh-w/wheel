package context

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey uint8

const (
	TransactionIDKey ctxKey = iota
)

func CtxWithTransactionID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, TransactionIDKey, id)
}

func TransactionIDFromCtx(ctx context.Context) uuid.UUID {
	v, _ := ctx.Value(TransactionIDKey).(uuid.UUID)

	return v
}
