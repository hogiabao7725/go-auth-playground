package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/config"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/health"
)

func main() {
	// load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// mux
	mux := http.NewServeMux()

	// register routes
	health.RegisterRoutes(mux)

	// server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	fmt.Printf("Server is up and running on PORT %s\n", cfg.Server.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}
