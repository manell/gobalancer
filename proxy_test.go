package gobalancer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGoBalancer(t *testing.T) {
	opts := &Options{}

	_, err := NewGoBalancer(opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoBalancerSimpleStrategy(t *testing.T) {
	response := "I'm the backend"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	}))
	defer backend.Close()

	opts := &Options{}

	gb, err := NewGoBalancer(opts)
	if err != nil {
		t.Fatal(err)
	}

	sb, err := NewSimpleBalancer(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	gb.UseStrategy(sb)

	req, err := http.NewRequest("GET", "http://var.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	gb.Proxy(w, req)

	if w.Body.String() != response {
		t.Log("Invalid response")
	}
}
