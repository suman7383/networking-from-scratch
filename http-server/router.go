package main

type Router struct {
	routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
	}
}

func (r *Router) HandleRoute(pattern string, handler HandlerFunc) {
	r.routes[pattern] = handler
}

func (r *Router) get(pattern string) (handler HandlerFunc, ok bool) {
	if h, ok := r.routes[pattern]; !ok {
		return nil, false
	} else {
		return h, true
	}
}
