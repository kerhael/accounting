package v1

type CreateCategoryRequest struct {
	Label string `json:"label"`
}

type CategoryResponse struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
