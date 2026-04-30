package auth

import "net/http"

type AuthRoutes struct {
	registerHL *RegisterHandler
	loginHL    *LoginHandler
	profileHL  *ProfileHandler
}

func NewAuthRoutes(registerHL *RegisterHandler, loginHL *LoginHandler, profileHL *ProfileHandler) *AuthRoutes {
	return &AuthRoutes{registerHL: registerHL, loginHL: loginHL, profileHL: profileHL}
}

func (ar *AuthRoutes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", ar.registerHL.HandleRegister)
	mux.HandleFunc("POST /auth/login", ar.loginHL.HandleLogin)
	mux.HandleFunc("GET /auth/profile", ar.profileHL.HandleProfile)
}
