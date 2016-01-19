package gobalancer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manell/gobalancer/strategy"
)

func TestGoBalancerSimpleStrategy(t *testing.T) {
	response := "I'm the backend"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	}))
	defer backend.Close()

	opts := &Options{}

	sb, err := strategy.NewSimpleBalancer(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	gb, err := NewGoBalancer(opts, sb)
	if err != nil {
		t.Fatal(err)
	}

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
