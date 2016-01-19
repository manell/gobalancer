package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var Key = []byte("There is no doubt that I'm a key")

// MockHandler is an HTPP handler that always return the value found in
// the test header
type MockHandler struct{}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(req.Header.Get("test")))
}

func createToken(t *testing.T, sub string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["sub"] = sub
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(Key)
	if err != nil {
		t.Fatal(err)
	}

	return tokenString
}

func writeNothing(token *jwt.Token, req *http.Request) error { return nil }

func writeSub(token *jwt.Token, req *http.Request) error {
	sub, _ := token.Claims["sub"]
	req.Header.Set("test", sub.(string))
	return nil
}

func TestJWT(t *testing.T) {
	tests := []struct {
		token    string
		function func(*jwt.Token, *http.Request) error
		status   int
		body     string
	}{
		{"Definitely I'm not a valid JWT", writeNothing, http.StatusUnauthorized, ""},
		{createToken(t, ""), writeNothing, http.StatusOK, ""},
		{createToken(t, ""), writeSub, http.StatusOK, ""},
		{createToken(t, "foo"), writeSub, http.StatusOK, "foo"},
	}

	for _, test := range tests {
		mid := &JWTAuth{
			KeyFunc: func(token *jwt.Token) (interface{}, error) {
				return Key, nil
			},
			ValidationFunction: test.function,
		}

		mock := &MockHandler{}
		h := mid.NewMiddleware(mock)

		req, err := http.NewRequest("GET", "http://var.com/foo", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "bearer "+test.token)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != test.status {
			t.Fatal("Invalid status code ", w.Code)
		}

		if w.Body.String() != test.body {
			t.Fatal("Invalid body ", w.Body)
		}
	}
}
