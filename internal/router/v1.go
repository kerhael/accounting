package router

import (
	"net/http"

	"github.com/kerhael/accounting/internal/auth"
	"github.com/kerhael/accounting/internal/handler"
	"github.com/kerhael/accounting/pkg/middleware"
)

func RegisterV1Routes(mux *http.ServeMux, h *handler.Handlers, rl *middleware.RateLimiter) {
	mux.HandleFunc("GET    /api/v1/health", h.V1.Health.Check)

	mux.Handle("GET    /api/v1/categories/", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Category.GetAllCategories)))
	mux.Handle("POST   /api/v1/categories/", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Category.PostCategory)))
	mux.Handle("GET    /api/v1/categories/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Category.GetCategoryById)))
	mux.Handle("DELETE /api/v1/categories/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Category.DeleteCategoryById)))

	mux.Handle("POST   /api/v1/outcomes/", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.PostOutcome)))
	mux.Handle("GET    /api/v1/outcomes/", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.GetAllOutcomes)))
	mux.Handle("GET    /api/v1/outcomes/sums-by-category", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.GetOutcomesSum)))
	mux.Handle("GET    /api/v1/outcomes/total", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.GetOutcomesTotal)))
	mux.Handle("GET    /api/v1/outcomes/series-by-category", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.GetOutcomesSeries)))
	mux.Handle("GET    /api/v1/outcomes/series-total", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.GetOutcomesTotalSeries)))
	mux.Handle("GET    /api/v1/outcomes/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.GetOutcomeById)))
	mux.Handle("PATCH  /api/v1/outcomes/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.PatchOutcomeById)))
	mux.Handle("DELETE /api/v1/outcomes/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Outcomes.DeleteOutcomeById)))

	mux.Handle("POST   /api/v1/incomes/", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Incomes.PostIncome)))
	mux.Handle("GET    /api/v1/incomes/", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Incomes.GetAllIncomes)))
	mux.Handle("GET    /api/v1/incomes/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Incomes.GetIncomeById)))
	mux.Handle("PATCH  /api/v1/incomes/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Incomes.PatchIncomeById)))
	mux.Handle("DELETE /api/v1/incomes/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Incomes.DeleteIncomeById)))

	mux.Handle("POST   /api/v1/users/", rl.RateLimitMiddleware(http.HandlerFunc(h.V1.Users.PostUser)))
	mux.Handle("GET    /api/v1/users/me", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Users.GetMe)))
	mux.Handle("PATCH  /api/v1/users/{id}", auth.AuthMiddleware(h.JWT)(http.HandlerFunc(h.V1.Users.PatchUserById)))

	mux.Handle("POST   /api/v1/login/", rl.RateLimitMiddleware(http.HandlerFunc(h.V1.Auth.Login)))
}
