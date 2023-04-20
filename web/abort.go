package web

import (
	"errors"
	"net/http"
)

func Abort(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")

	var (
		webErr *WebError
		code   = http.StatusInternalServerError
	)

	if errors.As(err, &webErr) {
		code = webErr.Code
	}

	return Respond(w, code, errorResponse{Error: err.Error()}) // nolint: wrapcheck
}

type errorResponse struct {
	Error string `json:"error"`
}
