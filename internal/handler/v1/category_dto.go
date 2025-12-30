package v1

import "github.com/kerhael/accounting/internal/domain"

type CreateCategoryRequest struct {
	Label string `json:"label"`
}

type CategoryResponse domain.Category

type CategoriesResponse []domain.Category
