package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/service/mocks"
	"github.com/kerhael/accounting/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestUserHandler_PostUser_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "john@example.com",
		"password":  "password123",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Create", ctx, "John", "Doe", "john@example.com", "password123").
		Return(&domain.User{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var data UserResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.ID)
	assert.Equal(t, "John", data.FirstName)
	assert.Equal(t, "Doe", data.LastName)
	assert.Equal(t, "john@example.com", data.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_PostUser_InvalidJSON(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader([]byte(`{invalid}`)))
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "Create")
}

func TestUserHandler_PostUser_MissingFirstName(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"lastName": "Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "Create")
}

func TestUserHandler_PostUser_MissingLastName(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"email":     "john@example.com",
		"password":  "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "Create")
}

func TestUserHandler_PostUser_MissingEmail(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"password":  "password123",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "Create")
}

func TestUserHandler_PostUser_MissingPassword(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "john@example.com",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "Create")
}

func TestUserHandler_PostUser_ShortPassword(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "john@example.com",
		"password":  "short",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertNotCalled(t, "Create")
}

func TestUserHandler_PostUser_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "john@example.com",
		"password":  "password123",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	invalidErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("email already in use")}
	mockService.On("Create", ctx, "John", "Doe", "john@example.com", "password123").
		Return(nil, invalidErr)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockService.AssertExpectations(t)
}

func TestUserHandler_PostUser_ServiceError(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := NewUserHandler(mockService)

	input := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "john@example.com",
		"password":  "password123",
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Create", ctx, "John", "Doe", "john@example.com", "password123").
		Return(nil, errors.New("db connection failed"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostUser(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	mockService.AssertExpectations(t)
}

// Test RateLimiter.
// The first 5 requests from the same IP should reach the user handler; the 6th should get 429.
func TestUserHandler_PostUsersRoute_RateLimiter_BurstOf5(t *testing.T) {
	const burst = 5

	callCount := 0
	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	// Same parameters as main.go: NewRateLimiter(1, 5)
	rl := middleware.NewRateLimiter(1, burst)
	handler := rl.Middleware(userHandler)

	// All burst requests should reach the user handler
	for i := 1; i <= burst; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", nil)
		req.RemoteAddr = "10.10.10.1:9000"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("request %d within burst: expected 201, got %d", i, w.Code)
		}
	}

	if callCount != burst {
		t.Errorf("expected user handler called %d times within burst, got %d", burst, callCount)
	}

	// The (burst+1)th request must be rate limited
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", nil)
	req.RemoteAddr = "10.10.10.1:9000"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("request beyond burst: expected 429, got %d", w.Code)
	}

	// User handler must not have been called on the blocked request
	if callCount != burst {
		t.Errorf("user handler must not be called when rate limited; call count: %d", callCount)
	}
}

// Test RateLimiter.
// Two clients hitting POST /api/v1/users/ each get their own rate-limit bucket.
func TestUserHandler_PostUsersRoute_RateLimiter_DifferentClientsAreIndependent(t *testing.T) {
	rl := middleware.NewRateLimiter(rate.Limit(0.001), 1) // burst of 1
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	// Client A exhausts their bucket
	reqA1 := httptest.NewRequest(http.MethodPost, "/api/v1/users/", nil)
	reqA1.RemoteAddr = "client-a:1111"
	handler.ServeHTTP(httptest.NewRecorder(), reqA1)

	reqA2 := httptest.NewRequest(http.MethodPost, "/api/v1/users/", nil)
	reqA2.RemoteAddr = "client-a:1111"
	wA2 := httptest.NewRecorder()
	handler.ServeHTTP(wA2, reqA2)

	if wA2.Code != http.StatusTooManyRequests {
		t.Errorf("client A second request: expected 429, got %d", wA2.Code)
	}

	// Client B's first request must still pass
	reqB := httptest.NewRequest(http.MethodPost, "/api/v1/users/", nil)
	reqB.RemoteAddr = "client-b:2222"
	wB := httptest.NewRecorder()
	handler.ServeHTTP(wB, reqB)

	if wB.Code != http.StatusCreated {
		t.Errorf("client B first request: expected 201, got %d", wB.Code)
	}
}
