package slogrules

import "log/slog"

func example() {
	// --- Правило: строчная буква ---

	slog.Info("Starting server on port 8080") // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect")           // want `log message must start with a lowercase letter`
	slog.Info("starting server")              // OK

	// --- Правило: только английский ---

	slog.Info("запуск сервера")      // want `log message must contain only English characters`
	slog.Error("ошибка подключения") // want `log message must contain only English characters`
	slog.Info("connecting to db")    // OK

	// --- Правило: спецсимволы ---

	slog.Info("server started!")             // want `log message must not contain special characters or emoji`
	slog.Error("connection failed!!!")       // want `log message must not contain special characters or emoji`
	slog.Warn("something went wrong...")     // want `log message must not contain special characters or emoji`
	slog.Info("server started successfully") // OK

	// --- Правило: чувствительные данные ---

	password := "hunter2"
	apiKey := "abc123"
	token := "xyz789"

	slog.Info("user password: " + password) // want `log message may contain sensitive data`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data`
	slog.Info("token: " + token)            // want `log message may contain sensitive data`
	slog.Info("user authenticated")         // OK — просто литерал, без конкатенации
	slog.Info("token validated")            // OK — просто литерал, без конкатенации

	// --- Корректные сообщения — диагностик быть не должно ---

	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("retrying in 5 seconds")
	slog.Debug("request payload size 1024 bytes")
}
