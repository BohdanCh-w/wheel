package api

import (
	"net/http"

	"github.com/gorilla/mux"
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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mx.ServeHTTP(w, req)
}

func (r *Router) wrapMiddleware(handler Handler, mid ...Middleware) Handler {
	fullMid := make([]Middleware, 0, len(r.mid)+len(mid))

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
