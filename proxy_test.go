package gobalancer

import (
	"testing"
)

func TestProxy(t *testing.T) {
	opts := &Options{}

	_, err := NewGoBalancer(opts)
	if err != nil {
		t.Fatal(err)
	}

}
