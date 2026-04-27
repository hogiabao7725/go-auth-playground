package auth

import "net/http"

type AuthRoutes struct {
	registerUC *RegisterHandler
}

func NewAuthRoutes(registerUC *RegisterHandler) *AuthRoutes {
	return &AuthRoutes{registerUC: registerUC}
}

func (ar *AuthRoutes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", ar.registerUC.RegisterHandle)
}
