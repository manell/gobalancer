package gobalancer

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Options is an object that defines the load balancer configuration.
type Options struct {
}

// EndPoint is a destination address where the balancer can proxy requests.
type EndPoint struct {
	URL *url.URL
}

// Balancer is an interface that represents the ability to return an endpoint
// provided an HTTP request.
type Balancer interface {
	NextEndpoint(http.Request) (*EndPoint, error)
}

// Dispatcher proxies HTPP requests to the endpoints selected by the balancer
type Dispatcher struct {
	Balancer  Balancer
	transport *http.Transport
}

// ServeHTTP is an HTTP handler that proxies a request to an endpoint
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ep, _ := d.Balancer.NextEndpoint(*req)

	proxy := httputil.NewSingleHostReverseProxy(ep.URL)

	rProxy := &httputil.ReverseProxy{
		Director:  proxy.Director,
		Transport: d.transport,
	}

	rProxy.ServeHTTP(w, req)
}

// GoBalancer is an HTTP handler that takes incomning requests and sends to
// another server. The destination server is determined by an algorithm.
type GoBalancer struct {
	middlewares *MiddlewareChain
}

// NewGoBalancer returns a new instance of a GoBalancer. It returns an error if
// some configuration is missing.
func NewGoBalancer(opt *Options, balancer Balancer) (*GoBalancer, error) {
	tp := &http.Transport{
		MaxIdleConnsPerHost: 200,
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
			// KeepAlive: 30 * time.Second,
		}).Dial,
		// TLSHandshakeTimeout: 10 * time.Second,
	}

	dispatcher := &Dispatcher{
		Balancer:  balancer,
		transport: tp,
	}

	gb := &GoBalancer{
		middlewares: &MiddlewareChain{chain: dispatcher},
	}

	return gb, nil
}

// Use adds a new middleware into the middlewares chain
func (gb *GoBalancer) Use(middleware Middleware) {
	gb.middlewares.Add(middleware)
}

// Proxy is an HTTP handler that run the defined middlewares chain
func (gb *GoBalancer) Proxy(w http.ResponseWriter, req *http.Request) {
	gb.middlewares.Run(w, req)
}
