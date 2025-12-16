package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kerhael/accounting/internal/repository"
)

type HealthHandler struct {
	repo repository.HealthRepository
}

func NewHealthHandler(repo repository.HealthRepository) *HealthHandler {
	return &HealthHandler{repo: repo}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := map[string]string{
		"server": "ok",
	}

	if err := h.repo.Check(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"db": "ko",
		})
		return
	}
	res["db"] = "ok"

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
