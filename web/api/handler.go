package api

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type WebsocketHandler func(ctx context.Context, c *websocket.Conn) error
