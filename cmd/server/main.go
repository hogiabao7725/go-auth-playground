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
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/token"
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/login"
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/profile"
	registerUC "github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/register"
)

func main() {
	ctx := context.Background()

	// 1. load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuartion", "error", err)
		os.Exit(1)
	}

	// 2. initialize logger
	appLogger := logger.New(logger.Options{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Pretty: cfg.Logger.Pretty,
		Env:    cfg.Server.Env,
		Writer: os.Stderr,
	})

	// 3. connect to database
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

	// 4. initialize infrastructure
	bcrypt := crypt.NewBcrypt()
	idGen := identifier.NewUUID()
	jwtProvider := token.NewJWT(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	// 5. initialize repositories
	queries := sqlc.New(dbpool)
	userRepo := persistence.NewUserRepository(queries)

	// 6. initialize use cases
	registerUC := registerUC.NewInteractor(bcrypt, idGen, userRepo)
	loginUC := login.NewInteractor(bcrypt, userRepo, jwtProvider)
	profileUC := profile.NewInteractor(userRepo)

	// 7. initialize handlers
	registerHandler := auth.NewRegisterHandler(registerUC)
	loginHandler := auth.NewLoginHandler(loginUC)
	profileHandler := auth.NewProfileHandler(profileUC)

	// 8. initialize middleware
	authMW := middleware.NewAuthMiddleware(jwtProvider)
	loggerMW := middleware.NewLoggerMiddleware(appLogger)

	// 9. define policy functions
	// 9. Define policy functions (app-specific, dùng Middleware type chuẩn)
	public := func(h http.Handler) http.Handler { return h }
	protected := authMW.RequireAuth()

	// 10. setup router
	mux := http.NewServeMux()

	// 11. register routes with appropriate middleware
	healthRoutes := health.NewHealthRoutes()
	healthRoutes.RegisterRoutes(mux, public)

	authRoutes := auth.NewAuthRoutes(registerHandler, loginHandler, profileHandler)
	authRoutes.RegisterRoutes(mux, public, protected)

	// ==================== GLOBAL MIDDLEWARE ====================
	handler := loggerMW.Handler()(mux)

	// logger register route
	appLogger.Info("registered routes")
	appLogger.Info("route", slog.String("method", "GET"), slog.String("pattern", "/health"))
	appLogger.Info("route", slog.String("method", "GET"), slog.String("pattern", "/auth/profile"))
	appLogger.Info("route", slog.String("method", "POST"), slog.String("pattern", "/auth/register"))
	appLogger.Info("route", slog.String("method", "POST"), slog.String("pattern", "/auth/login"))

	// 11. start server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	appLogger.Info("server is running on", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
