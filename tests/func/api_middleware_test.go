package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request %s %s", r.Method, r.RequestURI)
		next(w, r)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Not authentication", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Success")
}

func TestUnauthorize(t *testing.T) {
	handler := LoggingMiddleware(AuthMiddleware(Handler))
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSuccess(t *testing.T) {
	handler := LoggingMiddleware(AuthMiddleware(Handler))
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
