package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	whctx "github.com/bohdanch-w/wheel/context"
	wherr "github.com/bohdanch-w/wheel/errors"
	whlogger "github.com/bohdanch-w/wheel/logger"
	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type PanicMid struct {
	Logger whlogger.Logger
}

func (mid *PanicMid) Wrap(h api.Handler) api.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		transactionID := whctx.TransactionIDFromCtx(ctx)

		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()

				mid.Logger.
					WithTransaction(transactionID).
					With("panic", r).
					Errorf("Request got fatal server error: %s", stack)

				_ = web.Abort(w, &web.WebError{
					Code: http.StatusInternalServerError,
					Err:  wherr.Error("fatal server error"),
				})
			}
		}()

		return h(ctx, w, r)
	}

	return f
}
