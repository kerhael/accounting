package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/service"
)

type OutcomeHandler struct {
	service service.OutcomeServiceInterface
}

func NewOutcomeHandler(service service.OutcomeServiceInterface) *OutcomeHandler {
	return &OutcomeHandler{service: service}
}

// Create an outcome
// @Summary      Create an outcome
// @Description Create a new outcome
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param        outcome  body      CreateOutcomeRequest  true  "Outcome payload"
// @Success      201       {object}   OutcomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/ [post]
func (h *OutcomeHandler) PostOutcome(w http.ResponseWriter, r *http.Request) {
	var req CreateOutcomeRequest

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
	if req.CategoryId == 0 {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}
	if req.CreatedAt.IsZero() {
		http.Error(w, "creation date is required", http.StatusBadRequest)
		return
	}

	outcome, err := h.service.Create(r.Context(), req.Name, req.Amount, req.CategoryId, &req.CreatedAt)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toOutcomeResponse(outcome)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get all outcomes
// @Summary      Get all outcomes
// @Description  Retrieve all outcomes with optional date filtering (defaults to current month if not provided)
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to first day of current month)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Success      200   {array}   OutcomeResponse
// @Failure      400   {object}  ErrorResponse  "Bad request error"
// @Failure      404   {object}  ErrorResponse  "Not found error"
// @Failure      500   {object}  ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/ [get]
func (h *OutcomeHandler) GetAllOutcomes(w http.ResponseWriter, r *http.Request) {
	var from, to *time.Time
	var categoryId int

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

	categoryIdStr := r.URL.Query().Get("categoryId")
	if categoryIdStr != "" {
		categoryIdInt, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "invalid category", http.StatusBadRequest)
			return
		}
		categoryId = categoryIdInt
	}

	// If no dates provided, default to current month
	if from == nil && to == nil {
		now := time.Now()
		from = &time.Time{}
		*from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to = &now
	}

	outcomes, err := h.service.GetAll(r.Context(), from, to, categoryId)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidDateError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toOutcomesResponse(outcomes)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get an outcome
// @Summary      Get an outcome
// @Description Retrieve an outcome by id
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param 		id path int true "Outcome ID"
// @Success      200       {object}   OutcomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/{id} [get]
func (h *OutcomeHandler) GetOutcomeById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	outcome, err := h.service.GetById(r.Context(), id)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		if error, ok := errors.AsType[*domain.EntityNotFoundError](err); ok {
			http.Error(w, error.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toOutcomeResponse(outcome)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Update an outcome
// @Summary      Update an outcome
// @Description  Update an outcome
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param 		 id path int true "Outcome ID"
// @Param        outcome  body      PatchOutcomeRequest  true  "Outcome payload"
// @Success      200       {object}   OutcomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/{id} [patch]
func (h *OutcomeHandler) PatchOutcome(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req PatchOutcomeRequest
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
		if reqAmount < 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}
		amount = reqAmount

	}

	categoryId := 0
	if req.CategoryId != nil {
		reqCategoryId := *req.CategoryId
		if reqCategoryId < 0 {
			http.Error(w, "invalid category ID", http.StatusBadRequest)
			return
		}
		categoryId = reqCategoryId
	}

	outcome, err := h.service.Patch(r.Context(), id, name, amount, categoryId, req.CreatedAt)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		if error, ok := errors.AsType[*domain.EntityNotFoundError](err); ok {
			http.Error(w, error.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toOutcomeResponse(outcome)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Delete an outcome
// @Summary      Delete an outcome
// @Description Delete an outcome by id
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param 		id path int true "Outcome ID"
// @Success      204       "No Content"
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/{id} [delete]
func (h *OutcomeHandler) DeleteOutcomeById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteById(r.Context(), id)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get sum of outcomes by category
// @Summary      Get sum of outcomes by category
// @Description Get the total amount of outcomes by category between dates (defaults to current month if not provided), optionally filtered by category
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to first day of current month)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Param        categoryId query int false "Category ID filter"
// @Success      200   {object}   SumOutcomeResponse
// @Failure      400   {object}   ErrorResponse  "Bad request error"
// @Failure      404   {object}   ErrorResponse  "Not found error"
// @Failure      500   {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/sums-by-category [get]
func (h *OutcomeHandler) GetOutcomesSum(w http.ResponseWriter, r *http.Request) {
	var from, to *time.Time
	var categoryId int

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
		firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		from = &firstDayOfMonth
		to = &now
	}

	categoryIdStr := r.URL.Query().Get("categoryId")
	if categoryIdStr != "" {
		categoryIdInt, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "invalid category", http.StatusBadRequest)
			return
		}
		categoryId = categoryIdInt
	}

	categorySums, err := h.service.GetSum(r.Context(), from, to, categoryId)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidDateError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(categorySums); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get total of outcomes
// @Summary      Get total of outcomes
// @Description Get the total amount of outcomes between dates (defaults to current month if not provided)
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to first day of current month)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Success      200   {object}   TotalOutcomeResponse
// @Failure      400   {object}   ErrorResponse  "Bad request error"
// @Failure      500   {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/total [get]
func (h *OutcomeHandler) GetOutcomesTotal(w http.ResponseWriter, r *http.Request) {
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
		firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		from = &firstDayOfMonth
		to = &now
	}

	total, err := h.service.GetTotal(r.Context(), from, to)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidDateError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := TotalOutcomeResponse{Total: total}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get monthly series of outcomes
// @Summary      Get monthly series of outcomes
// @Description Get the sum of outcomes by category for each month between dates (defaults to last 12 months if not provided)
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to 12 months ago)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Success      200   {array}   SeriesOutcomeResponse
// @Failure      400   {object}   ErrorResponse  "Bad request error"
// @Failure      500   {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/series-by-category [get]
func (h *OutcomeHandler) GetOutcomesSeries(w http.ResponseWriter, r *http.Request) {
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

	// If only one date or no dates provided, default to (last) 12 months
	if from == nil || to == nil {
		if from == nil && to == nil {
			now := time.Now()
			twelveMonthsAgo := now.AddDate(0, -12, 0)
			from = &twelveMonthsAgo
			to = &now
		} else if from == nil {
			twelveMonthsAgo := to.AddDate(0, -12, 0)
			from = &twelveMonthsAgo
		} else {
			twelveMonthsAfter := from.AddDate(0, 12, 0)
			to = &twelveMonthsAfter
		}
	}

	series, err := h.service.GetSeries(r.Context(), from, to)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidDateError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(series); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get monthly series of outcomes' total amount
// @Summary      Get monthly series of outcomes' total amount
// @Description Get the total sum of outcomes for each month between dates (defaults to last 12 months if not provided)
// @Tags         outcomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to 12 months ago)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Success      200   {array}   TotalSeriesOutcomeResponse
// @Failure      400   {object}   ErrorResponse  "Bad request error"
// @Failure      500   {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/outcomes/series-total [get]
func (h *OutcomeHandler) GetOutcomesTotalSeries(w http.ResponseWriter, r *http.Request) {
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

	// If only one date or no dates provided, default to (last) 12 months
	if from == nil || to == nil {
		if from == nil && to == nil {
			now := time.Now()
			twelveMonthsAgo := now.AddDate(0, -12, 0)
			from = &twelveMonthsAgo
			to = &now
		} else if from == nil {
			twelveMonthsAgo := to.AddDate(0, -12, 0)
			from = &twelveMonthsAgo
		} else {
			twelveMonthsAfter := from.AddDate(0, 12, 0)
			to = &twelveMonthsAfter
		}
	}

	series, err := h.service.GetTotalSeries(r.Context(), from, to)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidDateError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(series); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func toOutcomeResponse(outcome *domain.Outcome) OutcomeResponse {
	return OutcomeResponse{
		Name:       outcome.Name,
		CreatedAt:  outcome.CreatedAt,
		Amount:     outcome.Amount,
		CategoryId: outcome.CategoryId,
		ID:         outcome.ID,
	}
}

func toOutcomesResponse(outcomes []domain.Outcome) []OutcomeResponse {
	var outcomesResp []OutcomeResponse
	if len(outcomes) == 0 {
		return []OutcomeResponse{}
	}
	for _, i := range outcomes {
		outcomesResp = append(outcomesResp, toOutcomeResponse(&i))
	}
	return outcomesResp
}
