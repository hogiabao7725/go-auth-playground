package auth

import (
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/register"
)

type RegisterHandler struct {
	uc *register.Interactor
}

func NewRegisterHandler(uc *register.Interactor) *RegisterHandler {
	return &RegisterHandler{uc: uc}
}
