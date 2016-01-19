package strategy

import (
	"net/http"
	"testing"
)

func TestSimpleBalancer(t *testing.T) {
	addr := "http://foo.com/var"
	sb, err := NewSimpleBalancer(addr)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://var.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	ep, err := sb.NextEndpoint(*req)
	if err != nil {
		t.Fatal(err)
	}

	if ep.URL.String() != addr {
		t.Log("Invalid endpoint")
	}
}
