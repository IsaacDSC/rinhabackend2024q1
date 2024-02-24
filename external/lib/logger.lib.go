package lib

import (
	"log/slog"
	"os"
)

func NewLogger(minLevel slog.Level) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: minLevel,
	}))

	return logger
}
