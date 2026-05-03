package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/request"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/login"
)

type LoginHandler struct {
	uc     login.LoginUseCase
	logger *slog.Logger
}

func NewLoginHandler(uc login.LoginUseCase, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{uc: uc, logger: logger}
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
	result, err := h.uc.Login(r.Context(), login.Command{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.HandleError(w, r, h.logger, err)
		return
	}

	// Set refresh token cookie
	setRefreshCookie(w, result.RefreshToken, result.RefreshTTL)

	// Send response
	res := LoginResponse{
		AccessToken:  result.AccessToken,
		ExpiresIn:    result.ExpiresIn,
		RefreshToken: result.RefreshToken,
		User: userInfo{
			ID:    result.User.ID(),
			Name:  result.User.Name().String(),
			Email: result.User.Email().String(),
			Role:  result.User.Role().String(),
		},
	}
	response.Success(w, http.StatusOK, "login successful", res)
}
