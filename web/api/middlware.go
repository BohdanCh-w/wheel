package api

type Middleware interface {
	Wrap(h Handler) Handler
}
