// example of use high order function to create middleware of api
package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func RateLimiterMiddleware(limit int) func(http.HandlerFunc) http.HandlerFunc {
	var countRequest int
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			print(limit)
			if countRequest >= limit {
				http.Error(w, "Too many request", http.StatusTooManyRequests)
				return
			}
			countRequest++
			next(w, r)
		}
	}
}

func CacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	var lastCall *time.Time
	return func(w http.ResponseWriter, r *http.Request) {
		if lastCall != nil && time.Now().Sub(*lastCall) <= 2*time.Second {
			fmt.Fprintf(w, "Cached")
			return
		}
		currentTime := time.Now()
		lastCall = &currentTime
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

func TestRateLimiterMiddleware(t *testing.T) {
	handler := LoggingMiddleware(RateLimiterMiddleware(1)(AuthMiddleware(Handler)))

	// First request should succeed
	req1 := httptest.NewRequest("GET", "/test", nil)
	rr1 := httptest.NewRecorder()
	req1.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	// Second request should fail due to rate limit
	req2 := httptest.NewRequest("GET", "/test", nil)
	rr2 := httptest.NewRecorder()
	req2.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rr2.Code)
}

func TestCacheMiddleware(t *testing.T) {
	handler := CacheMiddleware(Handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Success", rr.Body.String())

	req2 := httptest.NewRequest("GET", "/test", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusOK, rr2.Code)
	assert.Equal(t, "Cached", rr2.Body.String())
}
