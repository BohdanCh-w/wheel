package middleware

import (
	"context"
	"errors"
	"net/http"

	whlogger "github.com/bohdanch-w/wheel/logger"
	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type ErrorMid struct {
	Logger whlogger.Logger
}

func (mid *ErrorMid) Wrap(h api.Handler) api.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		transactionID := web.TransactionIDFromCtx(ctx)

		err := h(ctx, w, r)
		if err == nil {
			return nil
		}

		var webErr *web.WebError

		if !errors.As(err, &webErr) {
			webErr.Err = err
		}

		mid.Logger.
			WithTransaction(transactionID).
			With("status", webErr.Status()).
			Warnf("Request failed: %s", webErr.Error())

		if webErr.Status() > 0 {
			return web.Abort(w, webErr) // nolint: wrapcheck
		}

		return err
	}

	return f
}
