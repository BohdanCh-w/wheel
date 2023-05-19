package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type ErrorMid struct{}

func (mid *ErrorMid) Wrap(h api.Handler) api.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := h(ctx, w, r)
		if err == nil {
			return nil
		}

		var webErr *web.WebError

		if !errors.As(err, &webErr) {
			webErr.Err = err
		}

		return web.Abort(w, webErr) // nolint: wrapcheck
	}

	return f
}
