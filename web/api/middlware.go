package api

type Middleware interface {
	Wrap(Handler) Handler
}
