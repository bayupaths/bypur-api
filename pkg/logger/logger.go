package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

var fileWriter *lumberjack.Logger

func SetupLogger(levelStr, formatStr, logDir string) {
	level := parseLevel(levelStr)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	logFilePath := filepath.Join(logDir, "app.log")
	fileWriter = &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 5,
		LocalTime:  true,
		Compress:   true,
	}
	multi := io.MultiWriter(os.Stdout, fileWriter)

	var handler slog.Handler
	if strings.ToLower(formatStr) == "text" {
		handler = slog.NewTextHandler(multi, opts)
	} else {
		handler = slog.NewJSONHandler(multi, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func Close() {
	if fileWriter != nil {
		_ = fileWriter.Close()
	}
}

func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
