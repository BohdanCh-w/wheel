package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func NewRouter(mid ...Middleware) *Router {
	return &Router{
		mx:  mux.NewRouter(),
		mid: mid,
	}
}

type Router struct {
	mx  *mux.Router
	mid []Middleware
}

func (r *Router) RegisterRoute(route *Route) {
	handler := r.wrapMiddleware(route.Handler, route.Mid...)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := handler(ctx, w, r); err != nil {
			return
		}
	}

	r.mx.Handle(route.Path, http.HandlerFunc(h)).Name(route.Name).Methods(route.Methods...)
}

func (r *Router) RegisterFileRoute(route *FileRoute) {
	handler := func(_ context.Context, w http.ResponseWriter, r *http.Request) error { // nolint: unparam
		http.FileServer(http.Dir(route.Directory)).ServeHTTP(w, r)

		return nil
	}

	handler = r.wrapMiddleware(handler, route.Mid...)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := handler(ctx, w, r); err != nil {
			return
		}
	}

	r.mx.PathPrefix(route.Path).Handler(http.HandlerFunc(h)).Name(route.Name).Methods(http.MethodGet)
}

func (r *Router) RegisterWebsocketRoute(route *WebsocketRoute) {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		c, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return fmt.Errorf("failed to upgrade to websocket connection: %w", err)
		}

		defer c.Close()

		return route.Handler(ctx, &WSRequest{Origin: r, Conn: c})
	}

	handler = r.wrapMiddleware(handler, route.Mid...)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := handler(ctx, w, r); err != nil {
			return
		}
	}

	r.mx.Handle(route.Path, http.HandlerFunc(h)).Name(route.Name).Methods(http.MethodGet)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mx.ServeHTTP(w, req)
}

func (r *Router) wrapMiddleware(handler Handler, mid ...Middleware) Handler {
	fullMid := make([]Middleware, len(r.mid), len(r.mid)+len(mid))

	copy(fullMid, r.mid)
	fullMid = append(fullMid, mid...)

	for i := len(fullMid) - 1; i >= 0; i-- {
		h := fullMid[i]
		if h != nil {
			handler = h.Wrap(handler)
		}
	}

	return handler
}
