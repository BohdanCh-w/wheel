package api

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

var _ io.WriteCloser = (*WSRequest)(nil)

type WSRequest struct {
	Origin *http.Request
	Conn   *websocket.Conn
}

func (r *WSRequest) Write(p []byte) (n int, err error) {
	return len(p), r.Conn.WriteMessage(websocket.TextMessage, p)
}

func (r *WSRequest) Close() error {
	return r.Conn.Close()
}
