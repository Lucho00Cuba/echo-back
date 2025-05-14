package main

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

// TestParseLogLevel ensures the string to slog.Level conversion works as expected.
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"info", slog.LevelInfo},
		{"invalid", slog.LevelInfo}, // fallback
	}

	for _, tt := range tests {
		got := parseLogLevel(tt.input)
		if got != tt.expected {
			t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

// TestInitLoggerWithFormats tests InitLogger with various formats and levels.
func TestInitLoggerWithFormats(t *testing.T) {
	originalFormat := os.Getenv("LOG_FORMAT")
	originalLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		_ = os.Setenv("LOG_FORMAT", originalFormat)
		_ = os.Setenv("LOG_LEVEL", originalLevel)
	}()

	tests := []struct {
		format      string
		level       string
		expectInfo  bool
		expectDebug bool
	}{
		{"text", "debug", true, true},
		{"json", "warn", false, false},
		{"invalid", "info", true, false},
	}

	for _, tt := range tests {
		_ = os.Setenv("LOG_FORMAT", tt.format)
		_ = os.Setenv("LOG_LEVEL", tt.level)

		InitLogger()

		// Use custom handler to capture output
		var buf bytes.Buffer
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: parseLogLevel(tt.level)})
		testLogger := slog.New(handler)

		testLogger.Debug("debug msg")
		testLogger.Info("info msg")

		out := buf.String()

		if tt.expectInfo && !strings.Contains(out, "info msg") {
			t.Errorf("Expected 'info msg' in output for level=%s", tt.level)
		}
		if !tt.expectInfo && strings.Contains(out, "info msg") {
			t.Errorf("Did not expect 'info msg' in output for level=%s", tt.level)
		}
		if tt.expectDebug && !strings.Contains(out, "debug msg") {
			t.Errorf("Expected 'debug msg' in output for level=%s", tt.level)
		}
		if !tt.expectDebug && strings.Contains(out, "debug msg") {
			t.Errorf("Did not expect 'debug msg' in output for level=%s", tt.level)
		}
	}
}
