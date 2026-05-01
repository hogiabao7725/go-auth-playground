package health

import (
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
)

type HealthRoutes struct{}

func NewHealthRoutes() *HealthRoutes {
	return &HealthRoutes{}
}

func (hr *HealthRoutes) RegisterRoutes(mux *http.ServeMux, public middleware.Middleware) {
	mux.Handle("GET /health", public(http.HandlerFunc(hr.healthCheckHandler)))
}

func (hr *HealthRoutes) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "OK", nil)
}
