package auth

import (
	"errors"
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/request"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/login"
)

type LoginHandler struct {
	uc login.LoginUseCase
}

func NewLoginHandler(uc login.LoginUseCase) *LoginHandler {
	return &LoginHandler{uc: uc}
}

func (h *LoginHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Bind and validate request
	if err := request.BindJSON(r, &req); err != nil {
		if valErr, ok := errors.AsType[*request.ValidationError](err); ok {
			response.Error(w, http.StatusUnprocessableEntity, "validation error", valErr.Fields)
			return
		}

		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Execute use case
	result, err := h.uc.Execute(r.Context(), login.Command{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		statusCode, msg := response.MapDomainErrorToHTTP(err)
		response.Error(w, statusCode, msg, nil)
		return
	}

	res := LoginResponse{
		User: UserInfo{
			ID:    result.ID(),
			Name:  result.Name().String(),
			Email: result.Email().String(),
			Role:  result.Role().String(),
		},
	}
	response.Success(w, http.StatusOK, "login successful", res)
}
