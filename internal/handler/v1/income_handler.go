package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/service"
)

type IncomeHandler struct {
	service service.IncomeServiceInterface
}

func NewIncomeHandler(service service.IncomeServiceInterface) *IncomeHandler {
	return &IncomeHandler{service: service}
}

// Create an income
// @Summary      Create an income
// @Description Create a new income
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param        income  body      CreateIncomeRequest  true  "Income payload"
// @Success      201       {object}   IncomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/incomes/ [post]
func (h *IncomeHandler) PostIncome(w http.ResponseWriter, r *http.Request) {
	var req CreateIncomeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "amount is required and must be positive", http.StatusBadRequest)
		return
	}
	if req.CreatedAt.IsZero() {
		http.Error(w, "creation date is required", http.StatusBadRequest)
		return
	}

	income, err := h.service.Create(r.Context(), req.Name, req.Amount, &req.CreatedAt)
	if err != nil {
		if _, ok := err.(*domain.InvalidEntityError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(income)
}

// Get all incomes
// @Summary      Get all incomes
// @Description  Retrieve all incomes with optional date filtering (defaults to current month if not provided)
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to first day of current month)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Success      200   {array}   IncomeResponse
// @Failure      400   {object}  ErrorResponse  "Bad request error"
// @Failure      404   {object}  ErrorResponse  "Not found error"
// @Failure      500   {object}  ErrorResponse  "Internal server error"
// @Router       /api/v1/incomes/ [get]
func (h *IncomeHandler) GetAllIncomes(w http.ResponseWriter, r *http.Request) {
	var from, to *time.Time

	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		parsedFrom, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			http.Error(w, "invalid 'from' date format, use ISO 8601 (RFC3339)", http.StatusBadRequest)
			return
		}
		from = &parsedFrom
	}

	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		parsedTo, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			http.Error(w, "invalid 'to' date format, use ISO 8601 (RFC3339)", http.StatusBadRequest)
			return
		}
		to = &parsedTo
	}

	// If no dates provided, default to current month
	if from == nil && to == nil {
		now := time.Now()
		from = &time.Time{}
		*from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to = &now
	}

	incomes, err := h.service.GetAll(r.Context(), from, to)
	if err != nil {
		if _, ok := err.(*domain.InvalidDateError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incomes)
}

// Get an income
// @Summary      Get an income
// @Description Retrieve an income by id
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param 		id path int true "Income ID"
// @Success      200       {object}   IncomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/incomes/{id} [get]
func (h *IncomeHandler) GetIncomeById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	income, err := h.service.GetById(r.Context(), id)
	if err != nil {
		if _, ok := err.(*domain.InvalidEntityError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, ok := err.(*domain.EntityNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(income)
}

// Update an income
// @Summary      Update an income
// @Description  Update an income
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param 		 id path int true "Income ID"
// @Param        income  body      PatchIncomeRequest  true  "Income payload"
// @Success      200       {object}   IncomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/incomes/{id} [patch]
func (h *IncomeHandler) PatchIncome(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req PatchIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := ""
	if req.Name != nil {
		cleanName := strings.TrimSpace(*req.Name)
		name = cleanName
	}

	amount := 0
	if req.Amount != nil {
		reqAmount := *req.Amount
		if reqAmount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}
		amount = reqAmount

	}

	income, err := h.service.Patch(r.Context(), id, name, amount, req.CreatedAt)
	if err != nil {
		if _, ok := err.(*domain.EntityNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(income)
}

// Delete an income
// @Summary      Delete an income
// @Description Delete an income by id
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param 		id path int true "Income ID"
// @Success      204       "No Content"
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/incomes/{id} [delete]
func (h *IncomeHandler) DeleteIncomeById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteById(r.Context(), id)
	if err != nil {
		if _, ok := err.(*domain.InvalidEntityError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
