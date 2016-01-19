package strategy

import (
	"net/url"
)

// EndPoint is a destination address where the balancer can proxy requests.
type EndPoint struct {
	URL *url.URL
}
