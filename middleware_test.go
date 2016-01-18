package gobalancer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockMiddleware it's a simple middleware that writes a message
type MockMiddleware struct {
	done http.Handler
}

func (m *MockMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("mock"))
}

func (m *MockMiddleware) NewMiddleware(h http.Handler) http.Handler {
	m.done = h
	return m
}

// MockMiddleware2 it's a middleware that writes a message and calls the next
// middleware in the chain
type MockMiddleware2 struct {
	Count int
	done  http.Handler
}

func (m *MockMiddleware2) NewMiddleware(h http.Handler) http.Handler {
	m.done = h
	return m
}

func (m *MockMiddleware2) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.Count += 1
	m.done.ServeHTTP(w, req)
}

func TestMiddleware(t *testing.T) {
	mock := &MockMiddleware{}

	mid := &MiddlewareChain{}
	mid.Add(mock)

	req, err := http.NewRequest("GET", "http://var.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	mid.Run(w, req)

	if w.Body.String() != "mock" {
		t.Log("Invalid response")
	}
}

func TestMiddlewareChain(t *testing.T) {
	mock1 := &MockMiddleware{}
	mock2 := &MockMiddleware2{}
	mock3 := &MockMiddleware2{}

	mid := &MiddlewareChain{}

	mid.Add(mock1)
	mid.Add(mock2)
	mid.Add(mock3)

	req, err := http.NewRequest("GET", "http://var.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	mid.Run(w, req)

	if w.Body.String() != "mock" {
		t.Log("Invalid response")
	}
	if mock2.Count != 1 {
		t.Log("Count should be 1, but found %d", mock2.Count)
	}
	if mock3.Count != 1 {
		t.Log("Count should be 1, but found %d", mock3.Count)
	}
}
