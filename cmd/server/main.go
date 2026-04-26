package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/hogiabao7725/go-auth-playground/internal/config"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/health"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/logger"
)

func main() {
	// load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuartion", "error", err)
		os.Exit(1)
	}

	// logger
	appLogger := logger.New(logger.Options{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Pretty: cfg.Logger.Pretty,
		Env:    cfg.Server.Env,
		Writer: os.Stderr,
	})

	// mux
	mux := http.NewServeMux()

	// register routes
	health.RegisterRoutes(mux)

	// logger register route
	appLogger.Info("registered routes")
	appLogger.Info("route", slog.String("method", "GET"), slog.String("pattern", "/health"))

	// middleware
	lm := middleware.NewLoggerMiddleware(appLogger)
	handler := lm.Handler(mux)

	// server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	// start server
	appLogger.Info("server is running on", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
