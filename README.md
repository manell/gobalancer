# gobalancer
Gobalancer is a golang pluginable load balancer. Gobalancer is compatible with the go standard library http.handler, so it can be used with any other libraries that are compatible too.

## Installation
Just run: 

```go get github.com/manell/gobalancer```

## Usage
The simplest way to get an example runngin is by using the standard library and a simple strategy defined in the package.

```go
import (
	"log"
	"net/http"

	"github.com/manell/gobalancer"
	"github.com/manell/gobalancer/strategy"
)

// This is the single endpoint where we are going to redirect the requests
const endpoint = "http://127.0.0.1:8081"

func main() {
	opts := &gobalancer.Options{}

  // We can use the already defined in the librabry strategy to run a basic example
	strategy, err := strategy.NewSimpleBalancer(endpoint)
	if err != nil {
		log.Fatal(err)
	}

  // Create a new balancer handler using the previous strategy
	balancer, err := gobalancer.NewGoBalancer(opts, strategy)
	if err != nil {
		log.Fatal(err)
	}

  // Set up the balancer as a simple http.handler
	http.HandleFunc("/", balancer.Proxy)

	log.Fatal(http.ListenAndServe(":8080", nil))
}


```

## Built as a handler
The gobalancer itself it is just a go ```http.handler```, so it not provides any support on how you should run or configure the HTTP server that is going to host the balancer. Instead of doing this, the package provides an object that implements the ```http.handler``` interface, and lets the user decide about what server should run. What this really means is that this balancer is fully compatible with the standard go library ```http```, so deciding with package to use is up to you.


## Middlewares
## Strategy
