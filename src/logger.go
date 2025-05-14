package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

// InitLogger initializes the global logger with the configured level and format.
// It supports text or JSON output based on the LOG_FORMAT env variable.
//
// Set LOG_LEVEL to "debug", "info", "warn", "error"
// Set LOG_FORMAT to "text" or "json"
func InitLogger() {
	level := parseLogLevel(getEnv("LOG_LEVEL", "info"))
	format := getEnv("LOG_FORMAT", "text")

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// parseLogLevel converts string levels into slog.Level
func parseLogLevel(val string) slog.Level {
	switch val {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
