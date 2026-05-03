package response

import (
	"errors"
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

func MapDomainErrorToHTTP(err error) (int, string) {
	switch {
	case errors.Is(err, user.ErrEmptyID),
		errors.Is(err, user.ErrEmptyName),
		errors.Is(err, user.ErrEmptyEmail),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrEmptyPassword),
		errors.Is(err, user.ErrWeakPassword),
		errors.Is(err, user.ErrInvalidRole):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, user.ErrEmailAlreadyExists):
		return http.StatusConflict, err.Error()

	case errors.Is(err, user.ErrUserNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, user.ErrInvalidCredentials),
		errors.Is(err, user.ErrTokenInvalid),
		errors.Is(err, user.ErrTokenExpired),
		errors.Is(err, user.ErrRefreshTokenNotFound):
		return http.StatusUnauthorized, err.Error()

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
