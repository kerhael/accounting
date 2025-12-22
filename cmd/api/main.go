package main

import (
	"net/http"

	_ "github.com/kerhael/accounting/docs"

	"github.com/kerhael/accounting/internal/config"
	"github.com/kerhael/accounting/internal/db"
	"github.com/kerhael/accounting/internal/handler"
	"github.com/kerhael/accounting/internal/router"
	"github.com/kerhael/accounting/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Accounting API
// @version         1.0
// @description     API de suivi de comptabilité
// @termsOfService  https://example.com/terms/

// @contact.name   API Support
// @contact.email  kerhael.me@gmail.com

// @BasePath  /api/v1/

// @schemes http
func main() {
	logr := logger.New()

	// configuration
	cfg, err := config.Load()
	if err != nil {
		logr.Error("config error", err)
	}

	// database
	dbPool, err := db.NewPostgresPool(cfg.Database)
	if err != nil {
		logr.Error("db error", err)
	}
	defer dbPool.Close()

	// register handlers
	handlers := handler.NewHandlers(dbPool)

	// mux server
	mux := http.NewServeMux()

	// register routes
	router.RegisterRoutes(mux, handlers)

	// swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	if err := http.ListenAndServe(":8080", mux); err != http.ErrServerClosed {
		logr.Error("server error:", err)
	}
}
