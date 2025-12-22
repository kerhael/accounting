package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/service"
)

type CategoryHandler struct {
	service service.CategoryServiceInterface
}

func NewCategoryHandler(service service.CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// Create a category
// @Summary      Create a category
// @Description Create a new category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      CreateCategoryRequest  true  "Category payload"
// @Success      201       {object}   CategoryResponse
// @Failure      400       {object}   domain.ErrorResponse  "Bad request error"
// @Failure      500       {object}   domain.ErrorResponse  "Internal server error"
// @Router       /api/v1/categories [post]
func (h *CategoryHandler) PostCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Label == "" {
		http.Error(w, "label is required", http.StatusBadRequest)
		return
	}

	category, err := h.service.Create(r.Context(), req.Label)
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
	json.NewEncoder(w).Encode(category)
}

// Get a category
// @Summary      Get a category
// @Description Retrieve a category by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param 		id path int true "Category ID"
// @Success      200       {object}   CategoryResponse
// @Failure      400       {object}   domain.ErrorResponse  "Bad request error"
// @Failure      404       {object}   domain.ErrorResponse  "Not found error"
// @Failure      500       {object}   domain.ErrorResponse  "Internal server error"
// @Router       /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ivalid id", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetById(r.Context(), id)
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}
