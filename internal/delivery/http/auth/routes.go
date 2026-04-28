package auth

import "net/http"

type AuthRoutes struct {
	registerHL *RegisterHandler
	loginHL    *LoginHandler
}

func NewAuthRoutes(registerHL *RegisterHandler, loginHL *LoginHandler) *AuthRoutes {
	return &AuthRoutes{registerHL: registerHL, loginHL: loginHL}
}

func (ar *AuthRoutes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", ar.registerHL.HandleRegister)
	mux.HandleFunc("POST /auth/login", ar.loginHL.HandleLogin)
}
