package logger

import (
	"log/slog"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	cases := []struct {
		level  string
		format string
	}{
		{"debug", "json"},
		{"info", "text"},
		{"warn", "json"},
		{"error", "text"},
		{"unknown", "unknown"},
	}

	for _, tc := range cases {
		t.Run(tc.level+"-"+tc.format, func(t *testing.T) {
			SetupLogger(tc.level, tc.format)
			if slog.Default() == nil {
				t.Fatal("expected default logger to be configured")
			}
		})
	}
}
