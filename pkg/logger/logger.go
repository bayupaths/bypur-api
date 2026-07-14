package logger

import (
	"log/slog"
	"os"
	"strings"
)

// SetupLogger mengonfigurasi slog default berdasarkan format dan level dari ENV
func SetupLogger(levelStr, formatStr string) {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var slogHandler slog.Handler
	if strings.ToLower(formatStr) == "text" {
		slogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(slogHandler))
}
