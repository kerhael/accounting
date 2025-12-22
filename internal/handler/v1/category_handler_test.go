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

func TestCategoryHandler_PostCategory_Success(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	input := map[string]string{"label": "Food"}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Create", ctx, "Food").Return(&domain.Category{
		ID:    1,
		Label: "Food",
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostCategory(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var data domain.Category
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.ID)
	assert.Equal(t, "Food", data.Label)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_PostCategory_InvalidLabel(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	input := map[string]string{"label": ""}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostCategory(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_PostCategory_MissingLabelField(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	input := map[string]interface{}{}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostCategory(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_PostCategory_InvalidJSON(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	body := []byte(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostCategory(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_PostCategory_ServiceError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	input := map[string]string{"label": "Travel"}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Create", ctx, "Travel").Return(nil, errors.New("db failure"))

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostCategory(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestCategoryHandler_PostCategory_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	input := map[string]string{"label": "InvalidCategory"}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("category already exists")}
	mockService.On("Create", ctx, "InvalidCategory").Return(nil, invalidEntityErr)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostCategory(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_GetCategoryById_Success(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	expectedCategory := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockService.On("GetById", ctx, 1).Return(expectedCategory, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Category
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.ID)
	assert.Equal(t, "Food", data.Label)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryById_InvalidId(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/categories/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_GetCategoryById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid category id")}
	mockService.On("GetById", ctx, -1).Return(nil, invalidEntityErr)

	req := httptest.NewRequest(http.MethodGet, "/categories/-1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "-1")
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryById_EntityNotFoundError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("category not found")}
	mockService.On("GetById", ctx, 999).Return(nil, entityNotFoundErr)

	req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryById_ServiceError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	serviceErr := errors.New("database connection failed")
	mockService.On("GetById", ctx, 1).Return(nil, serviceErr)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_DeleteCategoryById_Success(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	mockService.On("DeleteById", ctx, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_DeleteCategoryById_InvalidId(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/categories/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_DeleteCategoryById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("cannot delete category with existing transactions")}
	mockService.On("DeleteById", ctx, 1).Return(invalidEntityErr)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_DeleteCategoryById_ServiceError(t *testing.T) {
	mockService := new(mocks.CategoryService)
	handler := NewCategoryHandler(mockService)

	ctx := context.Background()
	serviceErr := errors.New("database connection failed")
	mockService.On("DeleteById", ctx, 1).Return(serviceErr)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
