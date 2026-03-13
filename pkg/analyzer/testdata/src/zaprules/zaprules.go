package zaprules

import "go.uber.org/zap"

func example() {
	logger := zap.NewNop()

	// --- Строчная буква ---
	logger.Info("Starting server")  // want `log message must start with a lowercase letter`
	logger.Info("starting server")  // OK

	// --- Только английский ---
	logger.Error("ошибка сервера")  // want `log message must contain only English characters`
	logger.Error("server error")    // OK

	// --- Спецсимволы ---
	logger.Warn("warning!")         // want `log message must not contain special characters or emoji`
	logger.Warn("something is off") // OK

	// --- Чувствительные данные ---
	secret := "s3cr3t"
	logger.Info("secret: " + secret) // want `log message may contain sensitive data`
	logger.Info("operation complete") // OK

	// --- SugaredLogger ---
	sugar := logger.Sugar()

	sugar.Infow("Starting handler") // want `log message must start with a lowercase letter`
	sugar.Infow("starting handler") // OK
}
