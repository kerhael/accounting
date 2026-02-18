package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

// okHandler is a minimal next handler that returns 200.
func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestRateLimiter_AllowsRequestsWithinBurst(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(10), 3) // burst of 3
	handler := rl.Middleware(http.HandlerFunc(okHandler))

	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d within burst: expected 200, got %d", i, w.Code)
		}
	}
}

func TestRateLimiter_BlocksRequestsOverBurst(t *testing.T) {
	// Very slow refill (0.001 req/s) ensures the bucket won't refill during the test.
	rl := NewRateLimiter(rate.Limit(0.001), 2)
	handler := rl.Middleware(http.HandlerFunc(okHandler))

	ip := "10.0.0.1:5000"

	// First 2 requests should pass (burst)
	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d within burst: expected 200, got %d", i, w.Code)
		}
	}

	// 3rd request should be rate limited
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.RemoteAddr = ip
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("request beyond burst: expected 429, got %d", w.Code)
	}
}

func TestRateLimiter_DifferentIPsHaveSeparateLimits(t *testing.T) {
	// burst of 1 — first request from each IP passes, second is blocked
	rl := NewRateLimiter(rate.Limit(0.001), 1)
	handler := rl.Middleware(http.HandlerFunc(okHandler))

	// Exhaust the limit for IP1
	req1 := httptest.NewRequest(http.MethodPost, "/", nil)
	req1.RemoteAddr = "1.1.1.1:1000"
	handler.ServeHTTP(httptest.NewRecorder(), req1)

	// Second request from IP1 — should be blocked
	req2 := httptest.NewRequest(http.MethodPost, "/", nil)
	req2.RemoteAddr = "1.1.1.1:1000"
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("second request from same IP: expected 429, got %d", w2.Code)
	}

	// First request from IP2 — should pass (fresh bucket)
	req3 := httptest.NewRequest(http.MethodPost, "/", nil)
	req3.RemoteAddr = "2.2.2.2:2000"
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("first request from different IP: expected 200, got %d", w3.Code)
	}
}

func TestRateLimiter_CallsNextHandlerWhenAllowed(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	rl := NewRateLimiter(rate.Limit(10), 5)
	handler := rl.Middleware(next)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called, but it was not")
	}
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 from next handler, got %d", w.Code)
	}
}

func TestRateLimiter_DoesNotCallNextHandlerWhenBlocked(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	rl := NewRateLimiter(rate.Limit(0.001), 1)
	handler := rl.Middleware(next)

	ip := "192.168.0.1:4567"

	// Exhaust the burst
	req1 := httptest.NewRequest(http.MethodPost, "/", nil)
	req1.RemoteAddr = ip
	handler.ServeHTTP(httptest.NewRecorder(), req1)

	// This request should be blocked — next must NOT be called
	called = false
	req2 := httptest.NewRequest(http.MethodPost, "/", nil)
	req2.RemoteAddr = ip
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req2)

	if called {
		t.Fatal("expected next handler NOT to be called when rate limited")
	}
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}
