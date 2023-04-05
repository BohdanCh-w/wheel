package api

type Route struct {
	Name    string
	Path    string
	Mid     []Middleware
	Methods []string
	Handler Handler
}
