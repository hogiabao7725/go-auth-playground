package middleware

import "net/http"

type contextKey string

const (
	UserIDKey contextKey = "auth_user_id"
	RoleKey   contextKey = "auth_role"

	// Header constants
	authHeaderKey    = "Authorization"
	authHeaderPrefix = "Bearer"
)

type Middleware func(http.Handler) http.Handler
