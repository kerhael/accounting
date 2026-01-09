package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/handler"
)

func RegisterV1Routes(mux *http.ServeMux, h *handler.Handlers) {
	mux.HandleFunc("GET /api/v1/health", h.V1.Health.Check)

	mux.HandleFunc("GET /api/v1/categories/", h.V1.Category.GetAllCategories)
	mux.HandleFunc("POST /api/v1/categories/", h.V1.Category.PostCategory)
	mux.HandleFunc("GET /api/v1/categories/{id}", h.V1.Category.GetCategoryById)
	mux.HandleFunc("DELETE /api/v1/categories/{id}", h.V1.Category.DeleteCategoryById)

	mux.HandleFunc("POST /api/v1/outcomes/", h.V1.Outcomes.PostOutcome)
	mux.HandleFunc("GET /api/v1/outcomes/", h.V1.Outcomes.GetAllOutcomes)
	mux.HandleFunc("GET /api/v1/outcomes/{id}", h.V1.Outcomes.GetOutcomeById)
	mux.HandleFunc("PATCH /api/v1/outcomes/{id}", h.V1.Outcomes.PatchOutcome)
	mux.HandleFunc("DELETE /api/v1/outcomes/{id}", h.V1.Outcomes.DeleteOutcomeById)
}
