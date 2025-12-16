package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux, h *handler.Handlers) {
	mux.Handle("GET /health", h.Health)
}
