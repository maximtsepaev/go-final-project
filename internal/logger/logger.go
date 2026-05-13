// Package logger предоставляет функцию для инициализации логгера.
package logger

import (
	"log/slog"
	"os"
)

// InitLogger инициализирует логгер, который пишет в консоль
func InitLogger() *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	return logger
}
