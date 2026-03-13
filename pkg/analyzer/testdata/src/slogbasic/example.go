package slogbasic

import "log/slog"

func example() {
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("retrying in 5 seconds")
	slog.Debug("request payload received")

	logger := slog.Default()
	logger.Info("handler registered")
}
