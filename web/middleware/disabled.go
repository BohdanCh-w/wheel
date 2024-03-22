package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type DisabledMid struct {
	Disabled        bool
	ResponseCode    int
	ResponseMessage string
}

func (mid *DisabledMid) Wrap(h api.Handler) api.Handler {
	var (
		responseCode    = http.StatusNotFound
		responseMessage json.RawMessage
	)

	if mid.ResponseCode != 0 {
		responseCode = mid.ResponseCode
	}

	if mid.ResponseMessage != "" {
		responseMessage = json.RawMessage(`"` + mid.ResponseMessage + `"`)
	}

	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if !mid.Disabled {
			return h(ctx, w, r)
		}

		return web.Respond(w, responseCode, responseMessage)
	}

	return f
}
