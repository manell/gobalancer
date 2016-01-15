package gobalancer

import (
	"net/http"
	"net/url"
)

// SimpleBalancer is an endpoint dispatcher that always returns the same endpoint
type SimpleBalancer struct {
	url *url.URL
}

// NewSimpleBalancer return an instance of a SimpleBalancer
func NewSimpleBalancer(addr string) (sb *SimpleBalancer, err error) {
	sb = &SimpleBalancer{}
	sb.url, err = url.Parse(addr)

	return
}

// NextEndpoint always returns the endpoint confifured during the instantation
func (sb *SimpleBalancer) NextEndpoint(req http.Request) (*EndPoint, error) {
	return &EndPoint{URL: sb.url}, nil
}
