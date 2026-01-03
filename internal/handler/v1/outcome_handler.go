package v1

import (
	"encoding/json"
	"net/http"

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
// @Failure      400       {object}   domain.ErrorResponse  "Bad request error"
// @Failure      500       {object}   domain.ErrorResponse  "Internal server error"
// @Router       /api/v1/categories/ [post]
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
		if _, ok := err.(*domain.InvalidEntityError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(outcome)
}
