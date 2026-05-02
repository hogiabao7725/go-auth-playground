package user

import (
	"errors"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user/vo"
)

var (
	// ID
	ErrEmptyID = errors.New("id cannot be empty")

	// Email
	ErrEmailAlreadyExists = errors.New("this email is already registered")

	// User
	ErrUserNotFound = errors.New("user not found")

	// Authentication
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Token
	ErrTokenInvalid         = errors.New("token is invalid or unauthorized")
	ErrTokenExpired         = errors.New("session expired, please login again")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

// Validation errors (Aliased from VO)
var (
	ErrEmptyName     = vo.ErrEmptyName
	ErrEmptyEmail    = vo.ErrEmptyEmail
	ErrInvalidEmail  = vo.ErrInvalidEmail
	ErrEmptyPassword = vo.ErrEmptyPassword
	ErrWeakPassword  = vo.ErrWeakPassword
	ErrInvalidRole   = vo.ErrInvalidRole
)
