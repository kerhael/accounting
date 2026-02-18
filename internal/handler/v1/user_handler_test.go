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
	"github.com/stretchr/testify/assert"
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
