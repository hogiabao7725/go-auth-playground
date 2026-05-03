package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func sendJSON(w http.ResponseWriter, code int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, code int, message string, data any) {
	sendJSON(w, code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, message string, errors any) {
	sendJSON(w, code, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func HandleError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	if err == nil {
		return
	}

	// Map error → status code + message
	statusCode, msg := MapDomainErrorToHTTP(err)

	level := slog.LevelError
	if statusCode >= 400 && statusCode < 500 {
		level = slog.LevelWarn
	}

	logger.LogAttrs(r.Context(), level, "request failed",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("error", err.Error()),
		slog.String("remote_addr", r.RemoteAddr),
		slog.Int("status", statusCode),
		slog.String("user_agent", r.UserAgent()),
	)

	Error(w, statusCode, msg, nil)
}
