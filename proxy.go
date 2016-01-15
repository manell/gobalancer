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

// GoBalancer is an HTTP handler that takes incomning requests and sends to
// another server. The destination server is determined by an algorithm.
type GoBalancer struct {
	Balancer  Balancer
	transport *http.Transport
}

// EndPoint is a destination address where the balancer can proxy requests.
type EndPoint struct {
	URL *url.URL
}

// Balancer is an interface that represents the ability of returning and enpoint
// provided a HTTP request.
type Balancer interface {
	NextEndpoint(http.Request) (*EndPoint, error)
}

// NewGoBalancer returns a new instance of a GoBalancer. It returns an error if
// some configuration is missing.
func NewGoBalancer(opt *Options) (*GoBalancer, error) {
	tp := &http.Transport{
		MaxIdleConnsPerHost: 200,
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
			// KeepAlive: 30 * time.Second,
		}).Dial,
		// TLSHandshakeTimeout: 10 * time.Second,
	}

	gb := &GoBalancer{
		transport: tp,
	}

	return gb, nil
}

func (b *GoBalancer) Proxy(w http.ResponseWriter, req *http.Request) {
	ep, _ := b.Balancer.NextEndpoint(*req)

	proxy := httputil.NewSingleHostReverseProxy(ep.URL)

	rProxy := &httputil.ReverseProxy{
		Director:  proxy.Director,
		Transport: b.transport,
	}

	rProxy.ServeHTTP(w, req)
}
