package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

type MockHandlerlogs struct{}

func (m *MockHandlerlogs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

func TestLogs(t *testing.T) {
	// Redirect stdout to a pipe to properly check after the execution
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mid := NewLogs()

	mock := &MockHandlerlogs{}
	h := mid.NewMiddleware(mock)

	req, err := http.NewRequest("GET", "http://var.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	writ := httptest.NewRecorder()
	h.ServeHTTP(writ, req)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	// Check wether the final handler has been able to respond
	if writ.Body.String() != "OK" {
		t.Fatal("Invalid response")
	}

	request := regexp.MustCompile(`GET /foo`)
	// Cehck wether the log has been written
	if !request.MatchString(string(out)) {
		t.Fatal("Invalid log")
	}
}
