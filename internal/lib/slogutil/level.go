package slogutil

import "log/slog"

// LogLevel returns the slog.Level for the given target string.
func LogLevel(target string) slog.Level {
	switch target {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "warn", "WARN":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
