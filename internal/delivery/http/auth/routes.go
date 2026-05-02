package auth

import (
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
)

type AuthRoutes struct {
	registerHL *RegisterHandler
	loginHL    *LoginHandler
	profileHL  *ProfileHandler
	refreshHL  *RefreshHandler
}

func NewAuthRoutes(registerHL *RegisterHandler, loginHL *LoginHandler, profileHL *ProfileHandler, refreshHL *RefreshHandler) *AuthRoutes {
	return &AuthRoutes{registerHL: registerHL, loginHL: loginHL, profileHL: profileHL, refreshHL: refreshHL}
}

func (ar *AuthRoutes) RegisterRoutes(mux *http.ServeMux, public, protected middleware.Middleware) {
	mux.Handle("POST /auth/register", public(http.HandlerFunc(ar.registerHL.HandleRegister)))
	mux.Handle("POST /auth/login", public(http.HandlerFunc(ar.loginHL.HandleLogin)))
	mux.Handle("GET /auth/profile", protected(http.HandlerFunc(ar.profileHL.HandleProfile)))
	mux.Handle("POST /auth/refresh", public(http.HandlerFunc(ar.refreshHL.HandleRefresh)))
}
