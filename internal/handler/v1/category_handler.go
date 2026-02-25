package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/handler/utils"
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
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Label == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "label is required")
		return
	}

	category, err := h.service.Create(r.Context(), req.Label)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, toCategoryResponse(category))
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
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, toCategoriesResponse(categories))
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
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.service.GetById(r.Context(), id)
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

	utils.WriteJSON(w, http.StatusOK, toCategoryResponse(category))
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
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.service.DeleteById(r.Context(), id)
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

func toCategoryResponse(category *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:    category.ID,
		Label: category.Label,
	}
}

func toCategoriesResponse(categories []domain.Category) []CategoryResponse {
	var categoriesResp []CategoryResponse
	if len(categories) == 0 {
		return []CategoryResponse{}
	}
	for _, c := range categories {
		categoriesResp = append(categoriesResp, toCategoryResponse(&c))
	}
	return categoriesResp
}
