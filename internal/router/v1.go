package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/handler"
)

func RegisterV1Routes(mux *http.ServeMux, h *handler.Handlers) {
	mux.HandleFunc("GET /api/v1/health", h.V1.Health.Check)
	mux.HandleFunc("POST /api/v1/categories", h.V1.Category.PostCategory)
	mux.HandleFunc("/api/v1/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.V1.Category.GetCategoryById(w, r)
		case http.MethodDelete:
			h.V1.Category.DeleteCategoryById(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
