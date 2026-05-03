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

func TestIncomeHandler_PostIncome_Success(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	input := CreateIncomeRequest{
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: createdAt,
	}
	body, _ := json.Marshal(input)

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: &createdAt,
	}
	mockService.On("Create", ctx, "Salary", 300000, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	}), 123).Return(expectedIncome, nil)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var data domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.ID)
	assert.Equal(t, "Salary", data.Name)
	assert.Equal(t, 300000, data.Amount)
	assert.NotNil(t, data.CreatedAt)
	assert.True(t, data.CreatedAt.Equal(createdAt))

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PostIncome_NoAuthContext(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", nil)

	w := httptest.NewRecorder()
	handler.PostIncome(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "PostIncome")
}

func TestIncomeHandler_PostIncome_InvalidJSON(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader([]byte("invalid json")))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIncomeHandler_PostIncome_MissingName(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	input := CreateIncomeRequest{
		Amount:    300000,
		CreatedAt: time.Now(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "name is required")
}

func TestIncomeHandler_PostIncome_InvalidAmount(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	input := CreateIncomeRequest{
		Name:      "Salary",
		Amount:    0,
		CreatedAt: time.Now(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "amount is required and must be positive")
}

func TestIncomeHandler_PostIncome_NegativeAmount(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	input := CreateIncomeRequest{
		Name:      "Salary",
		Amount:    -100,
		CreatedAt: time.Now(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "amount is required and must be positive")
}

func TestIncomeHandler_PostIncome_ZeroCreatedAt(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	input := CreateIncomeRequest{
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: time.Time{},
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "creation date is required")
}

func TestIncomeHandler_PostIncome_ServiceError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	createdAt := time.Now()
	input := CreateIncomeRequest{
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: createdAt,
	}
	body, _ := json.Marshal(input)

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	mockService.On("Create", ctx, "Salary", 300000, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	}), 123).Return(nil, &domain.InvalidEntityError{UnderlyingCause: assert.AnError})

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_PostIncome_InternalError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	createdAt := time.Now()
	input := CreateIncomeRequest{
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: createdAt,
	}
	body, _ := json.Marshal(input)

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	mockService.On("Create", ctx, "Salary", 300000, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	}), 123).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.PostIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_Success(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	expectedIncomes := []domain.Income{
		{
			ID:        1,
			Name:      "Salary",
			Amount:    300000,
			CreatedAt: &time.Time{},
			UserId:    123,
		},
		{
			ID:        2,
			Name:      "Bonus",
			Amount:    50000,
			CreatedAt: &time.Time{},
			UserId:    123,
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 123, 20, 0).Return(expectedIncomes, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data PaginatedIncomesResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data.Data, 2)
	assert.Equal(t, expectedIncomes[0].ID, data.Data[0].ID)
	assert.Equal(t, expectedIncomes[0].Name, data.Data[0].Name)
	assert.Equal(t, expectedIncomes[1].ID, data.Data[1].ID)
	assert.Equal(t, expectedIncomes[1].Name, data.Data[1].Name)
	assert.Equal(t, 0, data.Pagination.Offset)
	assert.Equal(t, 20, data.Pagination.Limit)
	assert.Equal(t, 2, data.Pagination.Total)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetAllIncomes_NoAuthContext(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/", nil)

	w := httptest.NewRecorder()
	handler.GetAllIncomes(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetAllIncomes")
}

func TestIncomeHandler_GetAllIncomes_EmptyList(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedIncomes := []domain.Income{}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId, 20, 0).Return(expectedIncomes, 0, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data PaginatedIncomesResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data.Data, 0)
	assert.Empty(t, data.Data)
	assert.Equal(t, 0, data.Pagination.Total)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_WithDateFilters(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedIncomes := []domain.Income{
		{
			ID:        1,
			Name:      "Salary",
			Amount:    300000,
			CreatedAt: &time.Time{},
			UserId:    userId,
		},
	}
	mockService.On("GetAll", ctx, &from, &to, userId, 20, 0).Return(expectedIncomes, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data PaginatedIncomesResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data.Data, 1)
	assert.Equal(t, expectedIncomes[0].ID, data.Data[0].ID)
	assert.Equal(t, expectedIncomes[0].Name, data.Data[0].Name)
	assert.Equal(t, expectedIncomes[0].Amount, data.Data[0].Amount)
	assert.Equal(t, expectedIncomes[0].CreatedAt, data.Data[0].CreatedAt)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_WithPagination(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedIncomes := []domain.Income{
		{
			ID:        3,
			Name:      "Freelance",
			Amount:    120000,
			CreatedAt: &time.Time{},
			UserId:    userId,
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId, 10, 20).Return(expectedIncomes, 31, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?offset=20&limit=10", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data PaginatedIncomesResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data.Data, 1)
	assert.Equal(t, 20, data.Pagination.Offset)
	assert.Equal(t, 10, data.Pagination.Limit)
	assert.Equal(t, 31, data.Pagination.Total)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_InvalidOffset(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?offset=-1", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid offset")
	mockService.AssertNotCalled(t, "GetAll")
}

func TestIncomeHandler_GetAllIncomes_InvalidLimit(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?limit=101", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "limit must be less than or equal to 100")
	mockService.AssertNotCalled(t, "GetAll")
}

func TestIncomeHandler_GetAllIncomes_BadFromAndToDates(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId, 20, 0).Return([]domain.Income(nil), 0, invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "start date must be before end date")
}

func TestIncomeHandler_GetAllIncomes_InvalidFromDate(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?from=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'from' date format")
}

func TestIncomeHandler_GetAllIncomes_InvalidToDate(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?to=invalid-date", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid 'to' date format")
}

func TestIncomeHandler_GetAllIncomes_ServiceError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId, 20, 0).Return([]domain.Income(nil), 0, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/incomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetIncomeById_Success(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: &time.Time{},
		UserId:    userId,
	}
	mockService.On("GetById", ctx, 1, userId).Return(expectedIncome, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedIncome.ID, data.ID)
	assert.Equal(t, expectedIncome.Name, data.Name)
	assert.Equal(t, expectedIncome.Amount, data.Amount)
	assert.Equal(t, expectedIncome.CreatedAt, data.CreatedAt)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_GetIncomeById_NoAuthContext(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/1", nil)

	w := httptest.NewRecorder()
	handler.GetIncomeById(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "GetIncomeById")
}

func TestIncomeHandler_GetIncomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/invalid", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIncomeHandler_GetIncomeById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid income id")}
	mockService.On("GetById", ctx, -1, userId).Return(nil, invalidEntityErr)

	req := httptest.NewRequest(http.MethodGet, "/incomes/-1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "-1")
	w := httptest.NewRecorder()

	handler.GetIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetIncomeById_EntityNotFoundError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("income not found")}
	mockService.On("GetById", ctx, 999, userId).Return(nil, entityNotFoundErr)

	req := httptest.NewRequest(http.MethodGet, "/incomes/999", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetIncomeById_ServiceError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	serviceErr := errors.New("database connection failed")
	mockService.On("GetById", ctx, 1, userId).Return(nil, serviceErr)

	req := httptest.NewRequest(http.MethodGet, "/incomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_PatchIncomeById_Success_NameOnly(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      name,
		Amount:    300000,
		CreatedAt: &time.Time{},
		UserId:    userId,
	}
	mockService.On("PatchById", ctx, 1, name, 0, (*time.Time)(nil), userId).Return(expectedIncome, nil)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedIncome.ID, data.ID)
	assert.Equal(t, expectedIncome.Name, data.Name)
	assert.Equal(t, expectedIncome.Amount, data.Amount)
	assert.Equal(t, expectedIncome.CreatedAt, data.CreatedAt)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_PatchIncomeById_Success_AllFields(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	amount := 350000
	newCreatedAt := time.Now()
	userId := 123
	input := PatchIncomeByIdRequest{
		Name:      &name,
		Amount:    &amount,
		CreatedAt: &newCreatedAt,
	}
	body, _ := json.Marshal(input)

	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      name,
		Amount:    amount,
		CreatedAt: &newCreatedAt,
		UserId:    userId,
	}
	mockService.On("PatchById", ctx, 1, name, amount, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(newCreatedAt)
	}), userId).Return(expectedIncome, nil)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, expectedIncome.ID, data.ID)
	assert.Equal(t, expectedIncome.Name, data.Name)
	assert.Equal(t, expectedIncome.Amount, data.Amount)
	assert.True(t, data.CreatedAt.Equal(newCreatedAt))

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_PatchIncomeById_NoAuthContext(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", nil)

	w := httptest.NewRecorder()
	handler.PatchIncomeById(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "PatchIncomeById")
}

func TestIncomeHandler_PatchIncomeById_InvalidJSON(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader([]byte("invalid json")))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIncomeHandler_PatchIncomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/invalid", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid id")
}

func TestIncomeHandler_PatchIncomeById_NegativeAmount(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	amount := -100
	input := PatchIncomeByIdRequest{
		Amount: &amount,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "amount must be positive")
}

func TestIncomeHandler_PatchIncomeById_EntityNotFoundError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("income not found")}
	mockService.On("PatchById", ctx, 1, name, 0, (*time.Time)(nil), userId).Return(nil, entityNotFoundErr)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_PatchIncomeById_ServiceError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeByIdRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("PatchById", ctx, 1, name, 0, (*time.Time)(nil), userId).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_DeleteIncomeById_Success(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	userId := 123
	ctx := auth.ContextWithUserIDForTests(context.Background(), userId)
	mockService.On("DeleteById", ctx, 1, userId).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/incomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestOutcomeHandler_DeleteIncomeById_NoAuthContext(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/incomes/1", nil)

	w := httptest.NewRecorder()
	handler.DeleteIncomeById(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user not authenticated", response.Message)

	mockService.AssertNotCalled(t, "DeleteIncomeById")
}

func TestIncomeHandler_DeleteIncomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/incomes/invalid", nil)
	ctx := auth.ContextWithUserIDForTests(req.Context(), 123)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.DeleteIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid id")

	mockService.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything, 123)
}

func TestIncomeHandler_DeleteIncomeById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid income id")}
	mockService.On("DeleteById", ctx, 0, 123).Return(invalidEntityErr)

	req := httptest.NewRequest(http.MethodDelete, "/incomes/0", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	handler.DeleteIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_DeleteIncomeById_ServiceError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := auth.ContextWithUserIDForTests(context.Background(), 123)
	mockService.On("DeleteById", ctx, 1, 123).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/incomes/1", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
