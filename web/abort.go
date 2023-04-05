package web

import (
	"errors"
	"net/http"
)

func Abort(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")

	var webErr *WebError

	if errors.As(err, &webErr) {
		w.WriteHeader(webErr.Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, writeErr := w.Write([]byte(`{"error": "` + err.Error() + `"}`))

	return writeErr // nolint: wrapcheck
}
