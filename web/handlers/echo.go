package handlers

import (
	"context"
	"io"
	"net/http"

	whweb "github.com/bohdanch-w/wheel/web"
	whapi "github.com/bohdanch-w/wheel/web/api"
)

func Echo() whapi.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var req EchoRequest

		req.parse(r)

		return whweb.Respond(w, http.StatusOK, req)
	}
}

type EchoRequest struct {
	Body          string              `json:"body"`
	Sender        string              `json:"sender"`
	ContentLength int64               `json:"content_length"`
	Headers       map[string][]string `json:"headers"`
	Method        string              `json:"method"`
	URL           string              `json:"url"`
}

func (e *EchoRequest) parse(req *http.Request) {
	body, _ := io.ReadAll(req.Body)

	e.Body = string(body)
	e.Sender = req.RemoteAddr
	e.ContentLength = req.ContentLength
	e.Headers = req.Header
	e.Method = req.Method
	e.URL = req.URL.String()
}
