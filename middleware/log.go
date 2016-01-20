package middleware

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// Logs is a middleware thats logs info about the requests
type Logs struct {
	Log  *logrus.Logger
	done http.Handler
}

// NewLogs returns a new instance of Logs
func NewLogs() *Logs {
	log := logrus.New()

	return &Logs{Log: log}
}

// NewMiddleware registers the following middleware to be called after the execution
// of this
func (l *Logs) NewMiddleware(h http.Handler) http.Handler {
	l.done = h
	return l
}

// ServeHTTP implements the http.Handler interface. It logs the requests
func (l *Logs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	l.Log.WithFields(logrus.Fields{
		"src":     req.RemoteAddr,
		"request": fmt.Sprintf("%s %s %s", req.Method, req.URL.Path, req.Proto),
	}).Info()

	l.done.ServeHTTP(w, req)
}
