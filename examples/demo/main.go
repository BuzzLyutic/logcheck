// Этот файл демонстрирует работу линтера logcheck.
// Запуск: logcheck ./examples/demo/
package main

import "log/slog"

func main() {
	// Правило 1: заглавная буква - НАРУШЕНИЕ
	slog.Info("Starting server on port 8080")

	// Правило 2: не английский - НАРУШЕНИЕ
	slog.Error("ошибка подключения к базе данных")

	// Правило 3: спецсимволы - НАРУШЕНИЕ
	slog.Warn("something went wrong...")

	// Правило 4: чувствительные данные - НАРУШЕНИЕ
	token := "abc123"
	slog.Info("token: " + token)

	// Всё корректно - нарушений нет
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("retrying in 5 seconds")
	slog.Info("user authenticated successfully")
}
