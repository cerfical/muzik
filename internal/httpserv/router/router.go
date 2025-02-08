// Package router provides a simple and unified interface to perform HTTP routing.
package router

import (
	"fmt"
	"net/http"
	"strings"
)

// New constructs a new [Router].
func New() *Router {
	return &Router{mux: http.NewServeMux()}
}

// Router routes incoming requests to handlers as specified with [Router.Routes].
type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
}

// Middleware wraps [http.HandlerFunc] to perform additional actions before and/or after the original handler is called.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Endpoint describes the interface of a server endpoint defined by some URI path.
type Endpoint struct {
	// Method is the request method implemented by the endpoint.
	Method string

	// Handler is a function that will be called when the endpoint is requested and other parameters (such as the request method) are matched.
	Handler http.HandlerFunc
}

// ServeHTTP implements [http.Handler].
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply middleware
	h := http.HandlerFunc(r.mux.ServeHTTP)
	for _, m := range r.middleware {
		h = m(h)
	}
	h.ServeHTTP(w, req)
}

// Routes defines routes to reach a set of endpoints along the specified path.
func (r *Router) Routes(path string, endpoints []Endpoint) *Router {
	if strings.HasSuffix(path, "/") {
		// Disable wildcard behavior of a trailing slash
		path += "{$}"
	}

	for _, endpoint := range endpoints {
		r.mux.HandleFunc(fmt.Sprintf("%s %s", endpoint.Method, path), endpoint.Handler)
	}
	return r
}

// Use applies a [Middleware].
func (r *Router) Use(m Middleware) *Router {
	r.middleware = append(r.middleware, m)
	return r
}
