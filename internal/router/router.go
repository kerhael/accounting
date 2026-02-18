package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/handler"
	"github.com/kerhael/accounting/pkg/middleware"
)

func RegisterRoutes(mux *http.ServeMux, h *handler.Handlers, rl *middleware.RateLimiter) {
	RegisterV1Routes(mux, h, rl)
}
