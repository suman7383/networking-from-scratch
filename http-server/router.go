package main

type Router map[string]HandlerFunc

func (r Router) AddRoute(pattern string, handler HandlerFunc) {
	r[pattern] = handler
}
