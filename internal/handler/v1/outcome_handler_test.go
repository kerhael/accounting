package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOutcomeHandler_PostOutcome_Success(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	input := CreateOutcomeRequest{
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  createdAt,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  &createdAt,
	}
	mockService.On("Create", ctx, "Restaurant", 1999, 1, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	})).Return(expectedOutcome, nil)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var data domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.ID)
	assert.Equal(t, "Restaurant", data.Name)
	assert.Equal(t, 1999, data.Amount)
	assert.Equal(t, 1, data.CategoryId)
	assert.NotNil(t, data.CreatedAt)
	assert.True(t, data.CreatedAt.Equal(createdAt))

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PostOutcome_InvalidJSON(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestOutcomeHandler_PostOutcome_MissingName(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	input := CreateOutcomeRequest{
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  time.Now(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "name is required")
}

func TestOutcomeHandler_PostOutcome_InvalidAmount(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	input := CreateOutcomeRequest{
		Name:       "Restaurant",
		Amount:     0,
		CategoryId: 1,
		CreatedAt:  time.Now(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "amount is required and must be positive")
}

func TestOutcomeHandler_PostOutcome_MissingCategoryId(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	input := CreateOutcomeRequest{
		Name:      "Restaurant",
		Amount:    1999,
		CreatedAt: time.Now(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "category is required")
}

func TestOutcomeHandler_PostOutcome_ZeroCreatedAt(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	input := CreateOutcomeRequest{
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  time.Time{},
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "creation date is required")
}

func TestOutcomeHandler_PostOutcome_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	createdAt := time.Now()
	input := CreateOutcomeRequest{
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  createdAt,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Create", ctx, "Restaurant", 1999, 1, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	})).Return(nil, &domain.InvalidEntityError{UnderlyingCause: assert.AnError})

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PostOutcome_InternalError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	createdAt := time.Now()
	input := CreateOutcomeRequest{
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  createdAt,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Create", ctx, "Restaurant", 1999, 1, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	})).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostOutcome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
