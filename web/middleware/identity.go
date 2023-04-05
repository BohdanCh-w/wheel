package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	whlogger "github.com/bohdanch-w/wheel/logger"
	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type IdentityMid struct {
	Logger whlogger.Logger
}

func (mid *IdentityMid) Wrap(h api.Handler) api.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var (
			id    = uuid.New()
			start = time.Now()
		)

		ctx = web.CtxWithTransactionID(ctx, id)

		mid.Logger.
			WithTransaction(id).
			With("method", r.Method).
			With("at", start.Format("02-Jan-2006 15:04:05.999")).
			Infof("Request received: %s", r.URL.String())

		return h(ctx, w, r)
	}

	return f
}
