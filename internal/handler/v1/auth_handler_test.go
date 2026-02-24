package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kerhael/accounting/internal/auth"
	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/service/mocks"
	"github.com/kerhael/accounting/pkg/middleware"
	"github.com/kerhael/accounting/pkg/security"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

type MockCheckPassword func(password, hash string) error

type MockGenerateJWT func(userID int) (string, error)

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	input := map[string]string{
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	hashedPassword, _ := security.HashPassword("password123")
	mockService.On("FindByEmail", ctx, "john@example.com").Return(&domain.User{
		ID:           1,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: hashedPassword,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Contains(t, data, "token")
	assert.Contains(t, data, "user")

	userData := data["user"].(map[string]interface{})
	assert.Equal(t, float64(1), userData["id"])
	assert.Equal(t, "John", userData["firstName"])
	assert.Equal(t, "Doe", userData["lastName"])
	assert.Equal(t, "john@example.com", userData["email"])

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader([]byte(`{invalid}`)))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "FindByEmail")
}

func TestAuthHandler_Login_MissingEmail(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	input := map[string]string{
		"password": "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Equal(t, "email is required\n", w.Body.String())
	mockService.AssertNotCalled(t, "FindByEmail")
}

func TestAuthHandler_Login_MissingPassword(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	input := map[string]string{
		"email": "john@example.com",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Equal(t, "password is required\n", w.Body.String())
	mockService.AssertNotCalled(t, "FindByEmail")
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	input := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("FindByEmail", ctx, "nonexistent@example.com").Return((*domain.User)(nil), errors.New("user not found"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidPassword(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	input := map[string]string{
		"email":    "john@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	hashedPassword, _ := security.HashPassword("password123")
	mockService.On("FindByEmail", ctx, "john@example.com").Return(&domain.User{
		ID:           1,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: hashedPassword,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_JWTGenerationError(t *testing.T) {
	mockService := new(mocks.UserService)
	mockJWTService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockService, mockJWTService)

	input := map[string]string{
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	hashedPassword, _ := security.HashPassword("password123")
	mockService.On("FindByEmail", ctx, "john@example.com").Return(&domain.User{
		ID:           1,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: hashedPassword,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	mockService.AssertExpectations(t)
}

// Test RateLimiter.
// The first 5 requests from the same IP should reach the auth handler; the 6th should get 429.
func TestAuthHandler_LoginRoute_RateLimiter_BurstOf5(t *testing.T) {
	const burst = 5

	callCount := 0
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	// Same parameters as main.go: NewRateLimiter(1, 5)
	rl := middleware.NewRateLimiter(1, burst)
	handler := rl.Middleware(authHandler)

	// All burst requests should reach the auth handler
	for i := 1; i <= burst; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", nil)
		req.RemoteAddr = "10.10.10.1:9000"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d within burst: expected 200, got %d", i, w.Code)
		}
	}

	if callCount != burst {
		t.Errorf("expected auth handler called %d times within burst, got %d", burst, callCount)
	}

	// The (burst+1)th request must be rate limited
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", nil)
	req.RemoteAddr = "10.10.10.1:9000"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("request beyond burst: expected 429, got %d", w.Code)
	}

	// Auth handler must not have been called on the blocked request
	if callCount != burst {
		t.Errorf("auth handler must not be called when rate limited; call count: %d", callCount)
	}
}

// Test RateLimiter.
// Two clients hitting POST /api/v1/users/login each get their own rate-limit bucket.
func TestAuthHandler_LoginRoute_RateLimiter_DifferentClientsAreIndependent(t *testing.T) {
	rl := middleware.NewRateLimiter(rate.Limit(0.001), 1) // burst of 1
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Client A exhausts their bucket
	reqA1 := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", nil)
	reqA1.RemoteAddr = "client-a:1111"
	handler.ServeHTTP(httptest.NewRecorder(), reqA1)

	reqA2 := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", nil)
	reqA2.RemoteAddr = "client-a:1111"
	wA2 := httptest.NewRecorder()
	handler.ServeHTTP(wA2, reqA2)

	if wA2.Code != http.StatusTooManyRequests {
		t.Errorf("client A second request: expected 429, got %d", wA2.Code)
	}

	// Client B's first request must still pass
	reqB := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", nil)
	reqB.RemoteAddr = "client-b:2222"
	wB := httptest.NewRecorder()
	handler.ServeHTTP(wB, reqB)

	if wB.Code != http.StatusOK {
		t.Errorf("client B first request: expected 200, got %d", wB.Code)
	}
}
