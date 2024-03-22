package middleware

import (
	"context"
	"net/http"

	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type APIKeyAuthMid struct {
	APIKeyFunc func(r *http.Request) string
	APIKey     string
}

func (mid *APIKeyAuthMid) Wrap(h api.Handler) api.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		actualAPIKey := mid.APIKeyFunc(r)

		if actualAPIKey == "" {
			return web.Respond(w, http.StatusUnauthorized, map[string]string{"error": "missing API key"})
		}

		if actualAPIKey != mid.APIKey {
			return web.Respond(w, http.StatusUnauthorized, map[string]string{"error": "invalid API key"})
		}

		return h(ctx, w, r)
	}

	return f
}
