// Этот файл проверяет методы slog с context-аргументом.
// Сообщение находится не в первом аргументе — линтер должен
// корректно определить его позицию.
package slogcontext

import (
	"context"
	"log/slog"
)

func contextMethods() {
	ctx := context.Background()
	logger := slog.Default()

	// InfoContext: аргумент 0 — ctx, аргумент 1 — msg.
	logger.InfoContext(ctx, "Starting handler")   // want `log message must start with a lowercase letter`
	logger.InfoContext(ctx, "starting handler")   // OK

	logger.WarnContext(ctx, "ошибка обработки")   // want `log message must contain only English characters`
	logger.WarnContext(ctx, "processing failed")  // OK

	logger.ErrorContext(ctx, "something broke!")   // want `log message must not contain special characters or emoji`
	logger.ErrorContext(ctx, "something broke")    // OK

	logger.DebugContext(ctx, "request received")   // OK

	// Пакетные функции с Context.
	slog.InfoContext(ctx, "Handler started")       // want `log message must start with a lowercase letter`
	slog.InfoContext(ctx, "handler started")       // OK
}
