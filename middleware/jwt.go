package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// JWTAuth is a middleware that authenticates requests using JWT
type JWTAuth struct {
	// KeyFunc must return the key to be used to validating the JWT signature
	KeyFunc func(*jwt.Token) (interface{}, error)
	// The function that can optionally check any claims of the parsed JWT
	// it also can modify the request adding some headers or query parameters
	ValidationFunction func(*jwt.Token, *http.Request) error
	// Required to handle the next middleware in the chain
	done http.Handler
}

// NewMiddleware registers the following middleware to be called after the execution
// of this
func (m *JWTAuth) NewMiddleware(h http.Handler) http.Handler {
	m.done = h
	return m
}

// ServeHTTP implements the http.Handler interface. It will parse and validate the
// token. Then it will call the following middleware in the chain
func (m *JWTAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token, err := jwt.ParseFromRequest(req, m.KeyFunc)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if m.ValidationFunction(token, req) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m.done.ServeHTTP(w, req)
}
