package logger

import (
	"log/slog"
	"os"
)

func Setup(lvl slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
