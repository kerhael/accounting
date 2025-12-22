package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux, h *handler.Handlers) {
	RegisterV1Routes(mux, h)
}
