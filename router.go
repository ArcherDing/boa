package boa

type Router struct {
	boa    *Boa
	routes map[string]*Route
}

type Route struct {
	method   string
	handlers []HandlerFunc
}

func NewRouter(b *Boa) *Router {
	r := new(Router)
	r.boa = b
	return r
}

func (r *Router) Match(method string, path string, c *Context) ([]HandlerFunc, string) {
	for pattern, route := range r.routes {
		if pattern == path && route.method == method {
			return route.handlers, pattern
		}
	}
	return nil, ""
}

func (r *Router) Add(method, pattern string, handlers []HandlerFunc) {
	route := &Route{
		method:   method,
		handlers: handlers,
	}

	if r.routes == nil {
		r.routes = make(map[string]*Route)
	}
	r.routes[pattern] = route
}
