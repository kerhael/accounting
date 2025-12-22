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

func TestCategoryHandler_ServeHTTP_Success(t *testing.T) {
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

func TestCategoryHandler_ServeHTTP_InvalidLabel(t *testing.T) {
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

func TestCategoryHandler_ServeHTTP_InvalidJSON(t *testing.T) {
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

func TestCategoryHandler_ServeHTTP_ServiceError(t *testing.T) {
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
