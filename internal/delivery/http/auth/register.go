package auth

import (
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/request"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	registerUC "github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/register"
)

type RegisterHandler struct {
	uc *registerUC.Interactor
}

func NewRegisterHandler(uc *registerUC.Interactor) *RegisterHandler {
	return &RegisterHandler{uc: uc}
}

func (h *RegisterHandler) RegisterHandle(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Bind and validate request
	if err := request.BindJSON(w, r, &req); err != nil {
		if valErr, ok := err.(*request.ValidationError); ok {
			response.Error(w, http.StatusUnprocessableEntity, "validation error", valErr.Fields)
			return
		}

		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Execute use case
	result, err := h.uc.Execute(r.Context(), registerUC.Command{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		statusCode, msg := response.MapDomainErrorToHTTP(err)
		response.Error(w, statusCode, msg, nil)
		return
	}
	response.Success(w, http.StatusCreated, "user registered successfully", result)
}
