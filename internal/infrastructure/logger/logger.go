package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Options struct {
	Level  slog.Level
	Format string // "json" or "text"
	Pretty bool
	Writer io.Writer
}

func New(opts Options) *slog.Logger {
	w := opts.Writer
	if w == nil {
		w = os.Stderr
	}

	var handler slog.Handler
	handlerOptions := &slog.HandlerOptions{
		Level: opts.Level,
	}

	switch opts.Format {
	case "json":
		handler = slog.NewJSONHandler(w, handlerOptions)
	default: // text
		if opts.Pretty {
			handler = tint.NewHandler(w, &tint.Options{
				Level:   opts.Level,
				NoColor: false,
			})
		} else {
			handler = slog.NewTextHandler(w, handlerOptions)
		}
	}

	return slog.New(handler)
}

func ParseLevel(level string) slog.Level {
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
