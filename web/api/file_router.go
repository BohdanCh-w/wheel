package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func NewFileRouter(mid ...Middleware) *FileRouter {
	return &FileRouter{
		mx:  mux.NewRouter(),
		mid: mid,
	}
}

type FileRouter struct {
	mx  *mux.Router
	mid []Middleware
}

type FileRoute struct {
	Name      string
	Path      string
	Mid       []Middleware
	Directory string
}

func (r *FileRouter) RegisterRoute(route *FileRoute) {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	r.mx.Handle(route.Path, http.HandlerFunc(h)).Name(route.Name).Methods(http.MethodGet)
}

func (r *FileRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mx.ServeHTTP(w, req)
}

func (r *FileRouter) wrapMiddleware(handler Handler, mid ...Middleware) Handler {
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
