package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func TestOutcomeHandler_GetAllOutcomes_Success(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
		},
		{
			ID:         2,
			Name:       "Groceries",
			Amount:     5000,
			CategoryId: 2,
			CreatedAt:  &time.Time{},
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return(expectedOutcomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, expectedOutcomes[0].ID, data[0].ID)
	assert.Equal(t, expectedOutcomes[0].Name, data[0].Name)
	assert.Equal(t, expectedOutcomes[1].ID, data[1].ID)
	assert.Equal(t, expectedOutcomes[1].Name, data[1].Name)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetAllOutcomes_EmptyList(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	expectedOutcomes := []domain.Outcome{}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return(expectedOutcomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 0)
	assert.Empty(t, data)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetAllOutcomes_WithDateFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
		},
	}
	mockService.On("GetAll", ctx, &from, &to, 0).Return(expectedOutcomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, expectedOutcomes[0].ID, data[0].ID)
	assert.Equal(t, expectedOutcomes[0].Name, data[0].Name)
	assert.Equal(t, expectedOutcomes[0].Amount, data[0].Amount)
	assert.Equal(t, expectedOutcomes[0].CategoryId, data[0].CategoryId)
	assert.Equal(t, expectedOutcomes[0].CreatedAt, data[0].CreatedAt)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetAllOutcomes_WithCategoryFilter(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	categoryId := 1

	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: categoryId,
			CreatedAt:  &time.Time{},
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), categoryId).Return(expectedOutcomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?categoryId=1", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, expectedOutcomes[0].ID, data[0].ID)
	assert.Equal(t, expectedOutcomes[0].Name, data[0].Name)
	assert.Equal(t, expectedOutcomes[0].Amount, data[0].Amount)
	assert.Equal(t, expectedOutcomes[0].CategoryId, data[0].CategoryId)
	assert.Equal(t, expectedOutcomes[0].CreatedAt, data[0].CreatedAt)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetAllOutcomes_BadFromAndToDates(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return([]domain.Outcome(nil), invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "start date must be before end date")
}

func TestOutcomeHandler_GetAllOutcomes_InvalidFromDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?from=invalid-date", nil)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'from' date format")
}

func TestOutcomeHandler_GetAllOutcomes_InvalidToDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?to=invalid-date", nil)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'to' date format")
}

func TestOutcomeHandler_GetAllOutcomes_CategoryNotFound(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid category")}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 1).Return([]domain.Outcome(nil), invalidEntityErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?categoryId=1", nil)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid category")
}

func TestOutcomeHandler_GetAllOutcomes_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	ctx := context.Background()
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return([]domain.Outcome(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
