package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/hogiabao7725/go-auth-playground/internal/config"
	"github.com/hogiabao7725/go-auth-playground/internal/database"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/auth"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/health"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/crypt"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/identifier"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/logger"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/persistence"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/persistence/sqlc"
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/login"
	registerUC "github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/register"
)

func main() {
	ctx := context.Background()

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

	// database
	dbpool, err := database.NewPostgresPool(ctx, database.PoolConfig{
		DSN:            cfg.DB.DSN(),
		MaxConns:       cfg.DB.MaxConns,
		MinConns:       cfg.DB.MinConns,
		ConnLifetime:   cfg.DB.ConnLifetime,
		ConnIdleTime:   cfg.DB.ConnIdleTime,
		ConnectTimeout: cfg.DB.ConnectTimeout,
	})
	if err != nil {
		appLogger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbpool.Close()

	appLogger.Info("connected to database")

	// mux
	mux := http.NewServeMux()

	// infrastructure
	bcrypt := crypt.NewBcrypt()
	idGen := identifier.NewUUID()

	// repositories
	queries := sqlc.New(dbpool)
	userRepo := persistence.NewUserRepository(queries)

	// use case
	registerUC := registerUC.NewInteractor(bcrypt, idGen, userRepo)
	loginUC := login.NewInteractor(bcrypt, userRepo)

	// handlers
	registerHandler := auth.NewRegisterHandler(registerUC)
	loginHandler := auth.NewLoginHandler(loginUC)

	// register routes
	authRoutes := auth.NewAuthRoutes(registerHandler, loginHandler)

	health.RegisterRoutes(mux)
	authRoutes.RegisterRoutes(mux)

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
