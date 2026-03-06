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

	"github.com/kerhael/accounting/internal/auth"
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

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  &createdAt,
	}
	mockService.On("Create", ctx, "Restaurant", 1999, 1, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	}), 123).Return(expectedOutcome, nil)

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

func TestOutcomeHandler_PostOutcome_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", nil)

	w := httptest.NewRecorder()
	handler.PostOutcome(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "PostOutcome")
}

func TestOutcomeHandler_PostOutcome_InvalidJSON(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/outcomes/", bytes.NewReader([]byte("invalid json")))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	mockService.On("Create", ctx, "Restaurant", 1999, 1, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	}), 123).Return(nil, &domain.InvalidEntityError{UnderlyingCause: assert.AnError})

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

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	mockService.On("Create", ctx, "Restaurant", 1999, 1, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	}), 123).Return(nil, assert.AnError)

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

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
		{
			ID:         2,
			Name:       "Groceries",
			Amount:     5000,
			CategoryId: 2,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(expectedOutcomes, nil)

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

func TestOutcomeHandler_GetAllOutcomes_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/", nil)

	w := httptest.NewRecorder()
	handler.GetAllOutcomes(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetAllOutcomes")
}

func TestOutcomeHandler_GetAllOutcomes_EmptyList(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedOutcomes := []domain.Outcome{}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(expectedOutcomes, nil)

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

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
	}
	mockService.On("GetAll", ctx, &from, &to, 0, userId).Return(expectedOutcomes, nil)

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

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	categoryId := 1

	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: categoryId,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), categoryId, userId).Return(expectedOutcomes, nil)

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

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return([]domain.Outcome(nil), invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
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
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
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

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid category")}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 1, userId).Return([]domain.Outcome(nil), invalidEntityErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/?categoryId=1", nil)
	req = req.WithContext(ctx)
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

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return([]domain.Outcome(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllOutcomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomeById_Success(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockService.On("GetById", ctx, 1, userId).Return(expectedOutcome, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutcome.ID, data.ID)
	assert.Equal(t, expectedOutcome.Name, data.Name)
	assert.Equal(t, expectedOutcome.Amount, data.Amount)
	assert.Equal(t, expectedOutcome.CategoryId, data.CategoryId)
	assert.Equal(t, expectedOutcome.CreatedAt, data.CreatedAt)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomeById_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/1", nil)

	w := httptest.NewRecorder()
	handler.GetOutcomeById(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetOutcomeById")
}

func TestOutcomeHandler_GetOutcomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/invalid", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestOutcomeHandler_GetOutcomeById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid outcome id")}
	mockService.On("GetById", ctx, -1, userId).Return(nil, invalidEntityErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/-1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "-1")
	w := httptest.NewRecorder()

	handler.GetOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomeById_EntityNotFoundError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("outcome not found")}
	mockService.On("GetById", ctx, 999, userId).Return(nil, entityNotFoundErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/999", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomeById_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	serviceErr := errors.New("database connection failed")
	mockService.On("GetById", ctx, 1, userId).Return(nil, serviceErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PatchOutcomeById_Success_NameOnly(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	name := "Restaurant"
	input := PatchOutcomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       name,
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockService.On("PatchById", ctx, 1, name, 0, 0, (*time.Time)(nil), userId).Return(expectedOutcome, nil)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutcome.ID, data.ID)
	assert.Equal(t, expectedOutcome.Name, data.Name)
	assert.Equal(t, expectedOutcome.Amount, data.Amount)
	assert.Equal(t, expectedOutcome.CategoryId, data.CategoryId)
	assert.Equal(t, expectedOutcome.CreatedAt, data.CreatedAt)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PatchOutcomeById_Success_AllFields(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	name := "Restaurant"
	amount := 2000
	categoryId := 2
	newCreatedAt := time.Now()
	input := PatchOutcomeByIdRequest{
		Name:       &name,
		Amount:     &amount,
		CategoryId: &categoryId,
		CreatedAt:  &newCreatedAt,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       name,
		Amount:     amount,
		CategoryId: categoryId,
		CreatedAt:  &newCreatedAt,
		UserId:     userId,
	}
	mockService.On("PatchById", ctx, 1, name, amount, categoryId, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(newCreatedAt)
	}), userId).Return(expectedOutcome, nil)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Outcome
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutcome.ID, data.ID)
	assert.Equal(t, expectedOutcome.Name, data.Name)
	assert.Equal(t, expectedOutcome.Amount, data.Amount)
	assert.Equal(t, expectedOutcome.CategoryId, data.CategoryId)
	assert.True(t, data.CreatedAt.Equal(newCreatedAt))

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PatchOutcomeById_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", nil)

	w := httptest.NewRecorder()
	handler.PatchOutcomeById(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "PatchOutcomeById")
}

func TestOutcomeHandler_PatchOutcomeById_InvalidJSON(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader([]byte("invalid json")))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestOutcomeHandler_PatchOutcomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	name := "Restaurant"
	input := PatchOutcomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/invalid", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid id")
}

func TestOutcomeHandler_PatchOutcomeById_NegativeAmount(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	amount := -100
	input := PatchOutcomeByIdRequest{
		Amount: &amount,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "amount must be positive")
}

func TestOutcomeHandler_PatchOutcomeById_NegativeCategoryId(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	categoryId := -1
	input := PatchOutcomeByIdRequest{
		CategoryId: &categoryId,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid category ID")
}

func TestOutcomeHandler_PatchOutcomeById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	name := "Restaurant"
	input := PatchOutcomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid category")}
	mockService.On("PatchById", ctx, 1, name, 0, 0, (*time.Time)(nil), userId).Return(nil, invalidEntityErr)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PatchOutcomeById_EntityNotFoundError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	name := "Restaurant"
	input := PatchOutcomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("outcome not found")}
	mockService.On("PatchById", ctx, 1, name, 0, 0, (*time.Time)(nil), userId).Return(nil, entityNotFoundErr)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PatchOutcomeById_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	name := "Restaurant"
	input := PatchOutcomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("PatchById", ctx, 1, name, 0, 0, (*time.Time)(nil), userId).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPatch, "/outcomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_DeleteOutcomeById_Success(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("DeleteById", ctx, 1, userId).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/outcomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_DeleteOutcome_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/outcomes/1", nil)

	w := httptest.NewRecorder()
	handler.DeleteOutcomeById(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "DeleteOutcomeById")
}

func TestOutcomeHandler_DeleteOutcomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/outcomes/invalid", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.DeleteOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid id")

	mockService.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything)
}

func TestOutcomeHandler_DeleteOutcomeById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid outcome id")}
	mockService.On("DeleteById", ctx, 0, userId).Return(invalidEntityErr)

	req := httptest.NewRequest(http.MethodDelete, "/outcomes/0", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	handler.DeleteOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_DeleteOutcomeById_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("DeleteById", ctx, 1, userId).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/outcomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteOutcomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSum_Success_NoFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	categorySums := []domain.CategorySum{
		{CategoryId: 1, Total: 3000},
		{CategoryId: 2, Total: 1500},
	}
	mockService.On("GetSum", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(categorySums, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data SumOutcomeResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, 1, data[0].CategoryId)
	assert.Equal(t, 3000, data[0].Total)
	assert.Equal(t, 2, data[1].CategoryId)
	assert.Equal(t, 1500, data[1].Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSum_Success_WithFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	categorySums := []domain.CategorySum{
		{CategoryId: 1, Total: 3000},
	}
	mockService.On("GetSum", ctx, &from, &to, 1, userId).Return(categorySums, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z&categoryId=1", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data SumOutcomeResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, 1, data[0].CategoryId)
	assert.Equal(t, 3000, data[0].Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSum_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category", nil)

	w := httptest.NewRecorder()
	handler.GetOutcomesSum(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetOutcomesSum")
}

func TestOutcomeHandler_GetOutcomesSum_DefaultCurrentMonth(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	categorySums := []domain.CategorySum{
		{CategoryId: 1, Total: 3000},
	}
	mockService.On("GetSum", ctx, mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		expected := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return t.Equal(expected)
	}), mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		diff := now.Sub(*t)
		return diff >= 0 && diff < time.Second
	}), 0, userId).Return(categorySums, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data SumOutcomeResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, 1, data[0].CategoryId)
	assert.Equal(t, 3000, data[0].Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSum_InvalidFromDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category?from=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'from' date format")

	mockService.AssertNotCalled(t, "GetSum", mock.Anything, mock.Anything, mock.Anything, mock.Anything, 123)
}

func TestOutcomeHandler_GetOutcomesSum_InvalidToDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category?to=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'to' date format")

	mockService.AssertNotCalled(t, "GetSum", mock.Anything, mock.Anything, mock.Anything, mock.Anything, 123)
}

func TestOutcomeHandler_GetOutcomesSum_InvalidCategory(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category?categoryId=invalid", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid category")

	mockService.AssertNotCalled(t, "GetSum", mock.Anything, mock.Anything, mock.Anything, mock.Anything, 123)
}

func TestOutcomeHandler_GetOutcomesSum_InvalidDateError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetSum", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(nil, invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "start date must be before end date")

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSum_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid category")}
	mockService.On("GetSum", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 1, userId).Return(nil, invalidEntityErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category?categoryId=1", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid category")

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSum_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("GetSum", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/sums-by-category", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSum(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSeries_Success_NoFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedSeries := []domain.MonthlySeries{
		{
			Month: "2025-07",
			Categories: map[int]int{
				1: 3000,
				2: 1500,
				3: 0, // All categories included even with 0
			},
		},
		{
			Month: "2025-08",
			Categories: map[int]int{
				1: 2500,
				2: 0,
				3: 500,
			},
		},
	}
	mockService.On("GetSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(expectedSeries, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.MonthlySeries
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, "2025-07", data[0].Month)
	assert.Equal(t, map[int]int{1: 3000, 2: 1500, 3: 0}, data[0].Categories)
	assert.Equal(t, "2025-08", data[1].Month)
	assert.Equal(t, map[int]int{1: 2500, 2: 0, 3: 500}, data[1].Categories)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSeries_Success_WithFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedSeries := []domain.MonthlySeries{
		{
			Month: "2025-01",
			Categories: map[int]int{
				1: 3000,
			},
		},
	}
	mockService.On("GetSeries", ctx, &from, &to, userId).Return(expectedSeries, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.MonthlySeries
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "2025-01", data[0].Month)
	assert.Equal(t, map[int]int{1: 3000}, data[0].Categories)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSeries_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category", nil)

	w := httptest.NewRecorder()
	handler.GetOutcomesSeries(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetOutcomesSeries")
}

func TestOutcomeHandler_GetOutcomesSeries_DefaultLast12Months(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedSeries := []domain.MonthlySeries{}
	mockService.On("GetSeries", ctx, mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		expected := now.AddDate(0, -12, 0)
		diff := expected.Sub(*t)
		return diff >= 0 && diff < time.Second
	}), mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		diff := now.Sub(*t)
		return diff >= 0 && diff < time.Second
	}), userId).Return(expectedSeries, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.MonthlySeries
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 0)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSeries_InvalidFromDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category?from=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'from' date format")

	mockService.AssertNotCalled(t, "GetSeries", mock.Anything, mock.Anything, mock.Anything)
}

func TestOutcomeHandler_GetOutcomesSeries_InvalidToDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category?to=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'to' date format")

	mockService.AssertNotCalled(t, "GetSeries", mock.Anything, mock.Anything, mock.Anything)
}

func TestOutcomeHandler_GetOutcomesSeries_InvalidDateError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(nil, invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "start date must be before end date")

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesSeries_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("GetSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-by-category", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotal_Success_NoFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedTotal := 4500
	mockService.On("GetTotal", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(expectedTotal, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data TotalOutcomeResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, data.Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotal_Success_WithFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedTotal := 3000
	mockService.On("GetTotal", ctx, &from, &to, userId).Return(expectedTotal, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data TotalOutcomeResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, data.Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotal_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total", nil)

	w := httptest.NewRecorder()
	handler.GetOutcomesTotal(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetOutcomesTotal")
}

func TestOutcomeHandler_GetOutcomesTotal_DefaultCurrentMonth(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedTotal := 2500
	mockService.On("GetTotal", ctx, mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		expected := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return t.Equal(expected)
	}), mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		diff := now.Sub(*t)
		return diff >= 0 && diff < time.Second
	}), userId).Return(expectedTotal, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data TotalOutcomeResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, data.Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotal_InvalidFromDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total?from=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'from' date format")

	mockService.AssertNotCalled(t, "GetTotal", mock.Anything, mock.Anything, mock.Anything)
}

func TestOutcomeHandler_GetOutcomesTotal_InvalidToDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total?to=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'to' date format")

	mockService.AssertNotCalled(t, "GetTotal", mock.Anything, mock.Anything, mock.Anything)
}

func TestOutcomeHandler_GetOutcomesTotal_InvalidDateError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetTotal", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(0, invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "start date must be before end date")

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotal_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("GetTotal", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(0, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/total", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_Success_NoFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedSeries := []domain.MonthlyTotalSeries{
		{
			Month: "2025-07",
			Total: 3000,
		},
		{
			Month: "2025-08",
			Total: 2500,
		},
	}
	mockService.On("GetTotalSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(expectedSeries, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.MonthlyTotalSeries
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, "2025-07", data[0].Month)
	assert.Equal(t, 3000, data[0].Total)
	assert.Equal(t, "2025-08", data[1].Month)
	assert.Equal(t, 2500, data[1].Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_Success_WithFilters(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedSeries := []domain.MonthlyTotalSeries{
		{
			Month: "2025-01",
			Total: 3000,
		},
	}
	mockService.On("GetTotalSeries", ctx, &from, &to, userId).Return(expectedSeries, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.MonthlyTotalSeries
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "2025-01", data[0].Month)
	assert.Equal(t, 3000, data[0].Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_NoAuthContext(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total", nil)

	w := httptest.NewRecorder()
	handler.GetOutcomesTotalSeries(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetOutcomesTotalSeries")
}

func TestOutcomeHandler_GetOutcomesTotalSeries_DefaultLast12Months(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedSeries := []domain.MonthlyTotalSeries{}
	mockService.On("GetTotalSeries", ctx, mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		expected := now.AddDate(0, -12, 0)
		diff := expected.Sub(*t)
		return diff >= 0 && diff < time.Second
	}), mock.MatchedBy(func(t *time.Time) bool {
		now := time.Now()
		diff := now.Sub(*t)
		return diff >= 0 && diff < time.Second
	}), userId).Return(expectedSeries, nil)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.MonthlyTotalSeries
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 0)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_InvalidFromDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total?from=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'from' date format")

	mockService.AssertNotCalled(t, "GetTotalSeries", mock.Anything, mock.Anything, mock.Anything, 123)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_InvalidToDate(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total?to=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'to' date format")

	mockService.AssertNotCalled(t, "GetTotalSeries", mock.Anything, mock.Anything, mock.Anything, 123)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_InvalidDateError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetTotalSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(nil, invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "start date must be before end date")

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetOutcomesTotalSeries_ServiceError(t *testing.T) {
	mockService := new(mocks.OutcomeService)
	handler := NewOutcomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("GetTotalSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/outcomes/series-total", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetOutcomesTotalSeries(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
