package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/handler"
	"github.com/kerhael/accounting/pkg/middleware"
)

func RegisterV1Routes(mux *http.ServeMux, h *handler.Handlers, rl *middleware.RateLimiter) {
	mux.HandleFunc("GET    /api/v1/health", h.V1.Health.Check)

	mux.HandleFunc("GET    /api/v1/categories/", h.V1.Category.GetAllCategories)
	mux.HandleFunc("POST   /api/v1/categories/", h.V1.Category.PostCategory)
	mux.HandleFunc("GET    /api/v1/categories/{id}", h.V1.Category.GetCategoryById)
	mux.HandleFunc("DELETE /api/v1/categories/{id}", h.V1.Category.DeleteCategoryById)

	mux.HandleFunc("POST   /api/v1/outcomes/", h.V1.Outcomes.PostOutcome)
	mux.HandleFunc("GET    /api/v1/outcomes/", h.V1.Outcomes.GetAllOutcomes)
	mux.HandleFunc("GET    /api/v1/outcomes/sums-by-category", h.V1.Outcomes.GetOutcomesSum)
	mux.HandleFunc("GET    /api/v1/outcomes/total", h.V1.Outcomes.GetOutcomesTotal)
	mux.HandleFunc("GET    /api/v1/outcomes/series-by-category", h.V1.Outcomes.GetOutcomesSeries)
	mux.HandleFunc("GET    /api/v1/outcomes/series-total", h.V1.Outcomes.GetOutcomesTotalSeries)
	mux.HandleFunc("GET    /api/v1/outcomes/{id}", h.V1.Outcomes.GetOutcomeById)
	mux.HandleFunc("PATCH  /api/v1/outcomes/{id}", h.V1.Outcomes.PatchOutcome)
	mux.HandleFunc("DELETE /api/v1/outcomes/{id}", h.V1.Outcomes.DeleteOutcomeById)

	mux.HandleFunc("POST   /api/v1/incomes/", h.V1.Incomes.PostIncome)
	mux.HandleFunc("GET    /api/v1/incomes/", h.V1.Incomes.GetAllIncomes)
	mux.HandleFunc("GET    /api/v1/incomes/{id}", h.V1.Incomes.GetIncomeById)
	mux.HandleFunc("PATCH  /api/v1/incomes/{id}", h.V1.Incomes.PatchIncome)
	mux.HandleFunc("DELETE /api/v1/incomes/{id}", h.V1.Incomes.DeleteIncomeById)

	mux.Handle("POST       /api/v1/users/", rl.Middleware(http.HandlerFunc(h.V1.Users.PostUser)))

	mux.Handle("POST       /api/v1/login/", rl.Middleware(http.HandlerFunc(h.V1.Auth.Login)))
}
