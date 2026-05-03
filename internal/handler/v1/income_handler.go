package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kerhael/accounting/internal/auth"
	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/handler/utils"
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
// @Failure      401   {object}       ErrorResponse  "Unauthorized error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /incomes/ [post]
func (h *IncomeHandler) PostIncome(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req CreateIncomeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Amount <= 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "amount is required and must be positive")
		return
	}
	if req.CreatedAt.IsZero() {
		utils.WriteJSONError(w, http.StatusBadRequest, "creation date is required")
		return
	}

	income, err := h.service.Create(r.Context(), req.Name, req.Amount, &req.CreatedAt, userId)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, toIncomeResponse(income))
}

// Get all incomes
// @Summary      Get all incomes
// @Description  Retrieve all incomes with optional date filtering (defaults to current month if not provided)
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param        from  query     string  false  "Start date filter (ISO 8601 format, defaults to first day of current month)"
// @Param        to    query     string  false  "End date filter (ISO 8601 format, defaults to now)"
// @Param        offset query    int     false  "Items offset (defaults to 0)"
// @Param        limit query     int     false  "Items limit (defaults to 20, max 100)"
// @Success      200   {object}  PaginatedIncomesResponse
// @Failure      400   {object}  ErrorResponse  "Bad request error"
// @Failure      401   {object}  ErrorResponse  "Unauthorized error"
// @Failure      404   {object}  ErrorResponse  "Not found error"
// @Failure      500   {object}  ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /incomes/ [get]
func (h *IncomeHandler) GetAllIncomes(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var from, to *time.Time
	offset := domain.DefaultOffset
	limit := domain.DefaultLimit

	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		parsedFrom, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid 'from' date format, use ISO 8601 (RFC3339)")
			return
		}
		from = &parsedFrom
	}

	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		parsedTo, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid 'to' date format, use ISO 8601 (RFC3339)")
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

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = parsedOffset
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if parsedLimit > domain.MaxLimit {
			utils.WriteJSONError(w, http.StatusBadRequest, "limit must be less than or equal to 100")
			return
		}
		limit = parsedLimit
	}

	incomes, total, err := h.service.GetAll(r.Context(), from, to, userId, limit, offset)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidDateError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, PaginatedIncomesResponse{
		Data: toIncomesResponse(incomes),
		Pagination: PaginationResponse{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
	})
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
// @Failure      401       {object}  ErrorResponse  "Unauthorized error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /incomes/{id} [get]
func (h *IncomeHandler) GetIncomeById(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	income, err := h.service.GetById(r.Context(), id, userId)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		if error, ok := errors.AsType[*domain.EntityNotFoundError](err); ok {
			utils.WriteJSONError(w, http.StatusNotFound, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, toIncomeResponse(income))
}

// Update an income
// @Summary      Update an income
// @Description  Update an income
// @Tags         incomes
// @Accept       json
// @Produce      json
// @Param 		 id path int true "Income ID"
// @Param        income  body      PatchIncomeByIdRequest  true  "Income payload"
// @Success      200       {object}   IncomeResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      401       {object}   ErrorResponse  "Unauthorized error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /incomes/{id} [patch]
func (h *IncomeHandler) PatchIncomeById(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req PatchIncomeByIdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
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
			utils.WriteJSONError(w, http.StatusBadRequest, "amount must be positive")
			return
		}
		amount = reqAmount

	}

	income, err := h.service.PatchById(r.Context(), id, name, amount, req.CreatedAt, userId)
	if err != nil {
		if error, ok := errors.AsType[*domain.EntityNotFoundError](err); ok {
			utils.WriteJSONError(w, http.StatusNotFound, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, toIncomeResponse(income))
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
// @Failure      401       {object}   ErrorResponse  "Unauthorized error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /incomes/{id} [delete]
func (h *IncomeHandler) DeleteIncomeById(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.service.DeleteById(r.Context(), id, userId)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toIncomeResponse(income *domain.Income) IncomeResponse {
	return IncomeResponse{
		Name:      income.Name,
		Amount:    income.Amount,
		CreatedAt: income.CreatedAt,
		ID:        income.ID,
	}
}

func toIncomesResponse(incomes []domain.Income) []IncomeResponse {
	var incomesResp []IncomeResponse
	if len(incomes) == 0 {
		return []IncomeResponse{}
	}
	for _, i := range incomes {
		incomesResp = append(incomesResp, toIncomeResponse(&i))
	}
	return incomesResp
}
