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

func (lm *LoggerMiddleware) Handler() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the ResponseWriter to capture the status code
			rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			level := levelFromStatusCode(rec.statusCode)

			lm.logger.LogAttrs(
				r.Context(),
				level,
				"request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.Int("status", rec.statusCode),
				slog.Duration("duration", duration),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
		})
	}
}

// levelFromStatusCode maps HTTP status codes to slog levels for logging.
func levelFromStatusCode(status int) slog.Level {
	switch {
	case status >= 500:
		return slog.LevelError
	case status >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// statusRecorder wraps http.ResponseWriter to capture the HTTP response status code.
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and delegates to the underlying ResponseWriter.
func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}
