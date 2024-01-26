package api

type Route struct {
	Name    string
	Path    string
	Mid     []Middleware
	Methods []string
	Handler Handler
}

type FileRoute struct {
	Name      string
	Path      string
	Mid       []Middleware
	Directory string
}

type WebsocketRoute struct {
	Name    string
	Path    string
	Mid     []Middleware
	Handler WebsocketHandler
}
