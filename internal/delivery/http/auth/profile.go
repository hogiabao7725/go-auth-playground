package auth

import (
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	profileUC "github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/profile"
)

type ProfileHandler struct {
	uc profileUC.ProfileUseCase
}

func NewProfileHandler(uc profileUC.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

func (h *ProfileHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	cmd := profileUC.Command{UserID: userID}
	profileInfo, err := h.uc.GetProfile(r.Context(), cmd)
	if err != nil {
		status, msg := response.MapDomainErrorToHTTP(err)
		response.Error(w, status, msg, nil)
		return
	}

	response.Success(w, http.StatusOK, "profile retrieved successfully", profileInfo)
}
