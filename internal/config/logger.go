package config

import (
	"log/slog"
	"os"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Level  LogLevel `json:"level"`
	Format string   `json:"format"` // "json" or "text"
}

// NewLogger creates a new structured logger based on configuration
func NewLogger(cfg LoggerConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case LogLevelDebug:
		level = slog.LevelDebug
	case LogLevelInfo:
		level = slog.LevelInfo
	case LogLevelWarn:
		level = slog.LevelWarn
	case LogLevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // Add source file and line number
	}

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// Default to JSON for structured logging
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// SetDefaultLogger sets the default slog logger
func SetDefaultLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}
