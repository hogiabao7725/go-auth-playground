package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *slog.Logger
}

func NewLoggerMiddleware(logger *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

func (lm *LoggerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture the status code
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sr, r)

		duration := time.Since(start)

		lm.logger.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", sr.status),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}

// statusRecorder wraps http.ResponseWriter to capture the HTTP response status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code and delegates to the underlying ResponseWriter.
func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}
