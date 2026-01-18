package main

type Router struct {
	// routes maps a path to a handler func for handling the request.
	routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
	}
}
