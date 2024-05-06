package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	LevelDebug string = "DEBUG"
	LevelInfo  string = "INFO"
	LevelWarn  string = "WARN"
	LevelError string = "ERROR"

	slogFields ctxKey = "slog_fields"
)

type ContextHandler struct {
	slog.Handler
}

func NewSlogger(level string) *slog.Logger {

	handler := ContextHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: Level(level),
	})}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// Handle adds contextual attributes to the Record before calling the underlying handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	var v []slog.Attr
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

func Level(level string) slog.Level {
	switch {
	case level == LevelDebug:
		return slog.LevelDebug
	case level == LevelInfo:
		return slog.LevelInfo
	case level == LevelWarn:
		return slog.LevelWarn
	case level == LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
