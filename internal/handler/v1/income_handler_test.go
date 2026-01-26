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

	ctx := context.Background()
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: &createdAt,
	}
	mockService.On("Create", ctx, "Salary", 300000, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	})).Return(expectedIncome, nil)

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

func TestIncomeHandler_PostIncome_InvalidJSON(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/incomes/", bytes.NewReader([]byte("invalid json")))
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

	ctx := context.Background()
	mockService.On("Create", ctx, "Salary", 300000, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	})).Return(nil, &domain.InvalidEntityError{UnderlyingCause: assert.AnError})

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

	ctx := context.Background()
	mockService.On("Create", ctx, "Salary", 300000, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(createdAt)
	})).Return(nil, assert.AnError)

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

	ctx := context.Background()
	expectedIncomes := []domain.Income{
		{
			ID:        1,
			Name:      "Salary",
			Amount:    300000,
			CreatedAt: &time.Time{},
		},
		{
			ID:        2,
			Name:      "Bonus",
			Amount:    50000,
			CreatedAt: &time.Time{},
		},
	}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return(expectedIncomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, expectedIncomes[0].ID, data[0].ID)
	assert.Equal(t, expectedIncomes[0].Name, data[0].Name)
	assert.Equal(t, expectedIncomes[1].ID, data[1].ID)
	assert.Equal(t, expectedIncomes[1].Name, data[1].Name)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_EmptyList(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := context.Background()
	expectedIncomes := []domain.Income{}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return(expectedIncomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 0)
	assert.Empty(t, data)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_WithDateFilters(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := context.Background()
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedIncomes := []domain.Income{
		{
			ID:        1,
			Name:      "Salary",
			Amount:    300000,
			CreatedAt: &time.Time{},
		},
	}
	mockService.On("GetAll", ctx, &from, &to).Return(expectedIncomes, nil)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllIncomes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data []domain.Income
	err := json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, expectedIncomes[0].ID, data[0].ID)
	assert.Equal(t, expectedIncomes[0].Name, data[0].Name)
	assert.Equal(t, expectedIncomes[0].Amount, data[0].Amount)
	assert.Equal(t, expectedIncomes[0].CreatedAt, data[0].CreatedAt)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_GetAllIncomes_BadFromAndToDates(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := context.Background()
	invalidDatesErr := &domain.InvalidDateError{UnderlyingCause: errors.New("start date must be before end date")}
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return([]domain.Income(nil), invalidDatesErr)

	req := httptest.NewRequest(http.MethodGet, "/incomes/?from=2026-01-01T00:00:00Z&to=2025-01-01T00:00:00Z", nil)
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

	ctx := context.Background()
	mockService.On("GetAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return([]domain.Income(nil), assert.AnError)

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

	ctx := context.Background()
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      "Salary",
		Amount:    300000,
		CreatedAt: &time.Time{},
	}
	mockService.On("GetById", ctx, 1).Return(expectedIncome, nil)

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

func TestIncomeHandler_GetIncomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/incomes/invalid", nil)
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

	ctx := context.Background()
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid income id")}
	mockService.On("GetById", ctx, -1).Return(nil, invalidEntityErr)

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

	ctx := context.Background()
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("income not found")}
	mockService.On("GetById", ctx, 999).Return(nil, entityNotFoundErr)

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

	ctx := context.Background()
	serviceErr := errors.New("database connection failed")
	mockService.On("GetById", ctx, 1).Return(nil, serviceErr)

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

func TestIncomeHandler_PatchIncome_Success_NameOnly(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      name,
		Amount:    300000,
		CreatedAt: &time.Time{},
	}
	mockService.On("Patch", ctx, 1, name, 0, (*time.Time)(nil)).Return(expectedIncome, nil)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

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

func TestIncomeHandler_PatchIncome_Success_AllFields(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	amount := 350000
	newCreatedAt := time.Now()
	input := PatchIncomeRequest{
		Name:      &name,
		Amount:    &amount,
		CreatedAt: &newCreatedAt,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	expectedIncome := &domain.Income{
		ID:        1,
		Name:      name,
		Amount:    amount,
		CreatedAt: &newCreatedAt,
	}
	mockService.On("Patch", ctx, 1, name, amount, mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(newCreatedAt)
	})).Return(expectedIncome, nil)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

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

func TestIncomeHandler_PatchIncome_InvalidJSON(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader([]byte("invalid json")))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIncomeHandler_PatchIncome_InvalidId(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/invalid", bytes.NewReader(body))
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid id")
}

func TestIncomeHandler_PatchIncome_NegativeAmount(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	amount := -100
	input := PatchIncomeRequest{
		Amount: &amount,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "amount must be positive")
}

func TestIncomeHandler_PatchIncome_EntityNotFoundError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	entityNotFoundErr := &domain.EntityNotFoundError{UnderlyingCause: errors.New("income not found")}
	mockService.On("Patch", ctx, 1, name, 0, (*time.Time)(nil)).Return(nil, entityNotFoundErr)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_PatchIncome_ServiceError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	name := "Updated Salary"
	input := PatchIncomeRequest{
		Name: &name,
	}
	body, _ := json.Marshal(input)

	ctx := context.Background()
	mockService.On("Patch", ctx, 1, name, 0, (*time.Time)(nil)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPatch, "/incomes/1", bytes.NewReader(body))
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.PatchIncome(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestIncomeHandler_DeleteIncomeById_Success(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := context.Background()
	mockService.On("DeleteById", ctx, 1).Return(nil)

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

func TestIncomeHandler_DeleteIncomeById_InvalidId(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/incomes/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.DeleteIncomeById(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "invalid id")

	mockService.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything)
}

func TestIncomeHandler_DeleteIncomeById_InvalidEntityError(t *testing.T) {
	mockService := new(mocks.IncomeService)
	handler := NewIncomeHandler(mockService)

	ctx := context.Background()
	invalidEntityErr := &domain.InvalidEntityError{UnderlyingCause: errors.New("invalid income id")}
	mockService.On("DeleteById", ctx, 0).Return(invalidEntityErr)

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

	ctx := context.Background()
	mockService.On("DeleteById", ctx, 1).Return(assert.AnError)

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
