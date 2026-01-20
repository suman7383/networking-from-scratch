package main

type Router map[string]HandlerFunc

func (r Router) AddRoute(pattern string, handler HandlerFunc) {
	r[pattern] = handler
}

func (r Router) Get(pattern string) (handler HandlerFunc, ok bool) {
	if h, ok := r[pattern]; !ok {
		return nil, false
	} else {
		return h, true
	}
}
