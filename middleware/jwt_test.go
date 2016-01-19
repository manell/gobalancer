package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var Key = []byte("There is no doubt that I'm a key")

type MockHandler struct{}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hi"))
}

func createToken(t *testing.T) string {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["sub"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(Key)
	if err != nil {
		t.Fatal(err)
	}

	return tokenString
}

func TestJWT(t *testing.T) {
	mid := &JWTAuth{
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return Key, nil
		},
		ValidationFunction: func(token *jwt.Token, req *http.Request) error {
			return nil
		},
	}

	mock := &MockHandler{}
	h := mid.NewMiddleware(mock)

	req, err := http.NewRequest("GET", "http://var.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatal("Invalid status code %i", w.Code)
	}
}
