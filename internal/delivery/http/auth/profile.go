package auth

import (
	"log/slog"
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	profileUC "github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/profile"
)

type ProfileHandler struct {
	uc     profileUC.ProfileUseCase
	logger *slog.Logger
}

func NewProfileHandler(uc profileUC.ProfileUseCase, logger *slog.Logger) *ProfileHandler {
	return &ProfileHandler{uc: uc, logger: logger}
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
		response.HandleError(w, r, h.logger, err)
		return
	}

	response.Success(w, http.StatusOK, "profile retrieved successfully", profileInfo)
}
