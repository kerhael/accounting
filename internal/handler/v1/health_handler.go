package v1

import (
	"encoding/json"
	"net/http"

	"github.com/kerhael/accounting/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// Health check
// @Summary      Health check
// @Description Check server and database connectivity
// @Tags         health
// @Produce      plain
// @Success      200 {string} string '{"db":"ok","server":"ok"}'
// @Failure      503 {string} string '{"db":"ko","server":"ok"}'
// @Router       /api/v1/health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := map[string]string{
		"server": "ok",
	}

	if err := h.service.Check(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		err := json.NewEncoder(w).Encode(map[string]string{
			"db": "ko",
		})
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}
	res["db"] = "ok"

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
