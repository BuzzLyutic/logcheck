package fixlowercase

import "log/slog"

func example() {
	slog.Info("Starting server on port 8080") // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect")           // want `log message must start with a lowercase letter`
	slog.Info("starting server")              // OK — исправление не нужно
	slog.Debug("Request received")            // want `log message must start with a lowercase letter`
}
