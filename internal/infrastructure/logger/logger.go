package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Options struct {
	Level  string
	Format string // "text" or "json"
	Pretty bool
	Env    string // "development" or "production"
	Writer io.Writer
}

func New(opts Options) *slog.Logger {
	w := opts.Writer
	if w == nil {
		w = os.Stderr
	}

	level := parseLevel(opts.Level)

	var handler slog.Handler
	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	if opts.Env == "production" {
		handler = slog.NewJSONHandler(w, handlerOptions)
	} else { // development
		switch opts.Format {
		case "json":
			handler = slog.NewJSONHandler(w, handlerOptions)
		default:
			if opts.Pretty {
				handler = tint.NewHandler(w, &tint.Options{
					Level:  level,
					NoColor: false,
				})
			} else {
				handler = slog.NewTextHandler(w, handlerOptions)
			}
		}
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
