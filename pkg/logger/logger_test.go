package logger

import (
	"log/slog"
	"os"
	"testing"
)

func TestSetupLoggerConfiguresDefaultLogger(t *testing.T) {
	cases := []struct {
		level  string
		format string
	}{
		{"debug", "json"},
		{"info", "text"},
		{"warn", "json"},
		{"error", "text"},
		{"unknown", "unknown"}, // harus fallback ke Info tanpa panic
	}

	// Gunakan direktori temp agar tidak mencemari direktori kerja saat test
	tmpDir := t.TempDir()

	for _, tc := range cases {
		t.Run(tc.level+"-"+tc.format, func(t *testing.T) {
			SetupLogger(tc.level, tc.format, tmpDir)
			if slog.Default() == nil {
				t.Fatal("expected default logger to be configured, got nil")
			}
		})
	}
}

func TestSetupLoggerCreatesLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	SetupLogger("info", "json", tmpDir)

	// Pastikan file handle dilepas sebelum t.TempDir() mencoba hapus direktori
	t.Cleanup(Close)

	// Tulis satu log entry agar file terbentuk
	slog.Info("test log entry")

	// Flush: lumberjack membuka file secara lazy, trigger dengan Close + reopen implisit
	Close()

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read log dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected log file to be created in log directory, but directory is empty")
	}
}

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"unknown", slog.LevelInfo}, // fallback
		{"", slog.LevelInfo},        // empty fallback
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := parseLevel(tc.input)
			if got != tc.want {
				t.Fatalf("parseLevel(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
