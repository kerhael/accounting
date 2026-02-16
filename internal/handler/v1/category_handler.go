package v1

import (
	"encoding/json"
	"errors"
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
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/categories/ [post]
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
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// Get all categories
// @Summary      Get all categories
// @Description Retrieve all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200       {array}   CategoryResponse
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/categories/ [get]
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// Get a category
// @Summary      Get a category
// @Description Retrieve a category by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param 		id path int true "Category ID"
// @Success      200       {object}   CategoryResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetById(r.Context(), id)
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
	json.NewEncoder(w).Encode(category)
}

// Delete a category
// @Summary      Delete a category
// @Description Delete a category by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param 		id path int true "Category ID"
// @Success      204       "No Content"
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
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
