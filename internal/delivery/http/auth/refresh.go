package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	"github.com/hogiabao7725/go-auth-playground/internal/usecase/auth/refresh"
)

type RefreshHandler struct {
	uc     refresh.RefreshUseCase
	logger *slog.Logger
}

func NewRefreshHandler(uc refresh.RefreshUseCase, logger *slog.Logger) *RefreshHandler {
	return &RefreshHandler{uc: uc, logger: logger}
}

func (h *RefreshHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "refresh token not found", nil)
		return
	}

	// call use case
	result, err := h.uc.Refresh(r.Context(), refresh.Command{RawRefreshToken: cookie.Value})
	if err != nil {
		clearRefreshCookie(w)

		response.HandleError(w, r, h.logger, err)
		return
	}

	// set new refresh token cookie
	setRefreshCookie(w, result.RefreshToken, result.RefreshTTL)

	response.Success(w, http.StatusOK, "refresh token successfully", RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		ExpiresIn:    result.ExpiresIn,
		RefreshToken: result.RefreshToken,
	})
}

// helpers functions
func setRefreshCookie(w http.ResponseWriter, rawToken string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	})
}

func clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // xóa ngay
	})
}
