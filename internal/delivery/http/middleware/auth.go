package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

type AuthMiddleware struct {
	tokenProvider user.TokenProvider
}

func NewAuthMiddleware(tokenProvider user.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{tokenProvider: tokenProvider}
}

func (am *AuthMiddleware) RequireAuth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extract Authorization header
			authHeader := r.Header.Get(authHeaderKey)
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "authorization header is required", nil)
				return
			}

			// 2. validate format: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], authHeaderPrefix) {
				response.Error(w, http.StatusUnauthorized, "authorization header format must be Bearer {token}", nil)
				return
			}

			// 3. Parse token
			claims, err := am.tokenProvider.ParseAccessToken(parts[1])
			if err != nil {
				if errors.Is(err, user.ErrTokenExpired) {
					response.Error(w, http.StatusUnauthorized, "token expired", nil)
					return
				}
				response.Error(w, http.StatusUnauthorized, "invalid or expired token", nil)
				return
			}

			// 4. Store user info in context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	val, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return val
}

func GetRole(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	val, ok := ctx.Value(RoleKey).(string)
	if !ok {
		return ""
	}
	return val
}

func RequireRole(allowedRoles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetRole(r.Context())
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Error(w, http.StatusForbidden, "you do not have permission to access this resource", nil)
		})
	}
}
