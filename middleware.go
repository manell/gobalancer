package gobalancer

import (
	"net/http"
)

// Middleware is an interface that represents the ability to create a middleware
type Middleware interface {
	NewMiddleware(http.Handler) http.Handler
}

// MiddlewareChain contains the configured middlewares chain
type MiddlewareChain struct {
	chain http.Handler
}

// Add adds a new middleware to the chain
func (m *MiddlewareChain) Add(middleware Middleware) {
	m.chain = middleware.NewMiddleware(m.chain)
}

// Run runs the middleware chain
func (m *MiddlewareChain) Run(w http.ResponseWriter, req *http.Request) {
	m.chain.ServeHTTP(w, req)
}
