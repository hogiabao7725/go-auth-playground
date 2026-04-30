package main

import (
	"net/http"

	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/middleware"
)

type MiddlewareStack func(http.HandlerFunc) http.Handler

// Public: policy — no authentication required
func Public(h http.HandlerFunc) http.Handler {
	return h
}

// WithAuth: policy — need valid JWT token
func WithAuth(mw *middleware.AuthMiddleware) MiddlewareStack {
	return func(h http.HandlerFunc) http.Handler {
		return mw.RequireAuth(h)
	}
}

// WithRole: policy — need valid JWT token + user must have at least one of the specified roles
func WithRole(mw *middleware.AuthMiddleware, roles ...string) MiddlewareStack {
	return func(h http.HandlerFunc) http.Handler {
		return mw.RequireAuth(
			middleware.RequireRole(roles...)(h),
		)
	}
}
