package main

import (
	"net/http"

	"github.com/kerhael/accounting/internal/config"
	"github.com/kerhael/accounting/internal/db"
	"github.com/kerhael/accounting/internal/handler"
	"github.com/kerhael/accounting/internal/router"
	"github.com/kerhael/accounting/pkg/logger"
)

func main() {
	logr := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logr.Error("config error", err)
	}

	dbPool, err := db.NewPostgresPool(cfg.Database)
	if err != nil {
		logr.Error("db error", err)
	}
	defer dbPool.Close()

	handlers := handler.NewHandlers(dbPool)

	mux := http.NewServeMux()

	router.RegisterRoutes(mux, handlers)

	if err := http.ListenAndServe(":8080", mux); err != http.ErrServerClosed {
		logr.Error("server error:", err)
	}
}
