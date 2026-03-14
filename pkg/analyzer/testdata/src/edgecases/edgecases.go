// Этот файл проверяет пограничные случаи для всех правил.
package edgecases

import "log/slog"

func edgeCases() {
	// === LOWERCASE ===

	// Строка начинается с цифры — это не uppercase, ОК.
	slog.Info("3 retries remaining")

	// Строка начинается с пунктуации — ОК.
	slog.Info("[server] ready")

	// Одна заглавная буква.
	slog.Info("A")  // want `log message must start with a lowercase letter`

	// Одна строчная буква — ОК.
	slog.Info("a")

	// Пустая строка — нечего проверять, ОК.
	slog.Info("")

	// Заглавная не-ASCII буква (немецкая Ü) — ловим.
	slog.Info("Über alles") // want `log message must start with a lowercase letter` `log message must contain only English characters`

	// === ENGLISH ===

	// Акцентированные латинские буквы (> MaxASCII) — ловим.
	slog.Info("café connection") // want `log message must contain only English characters`

	// Японские символы.
	slog.Info("サーバー起動")  // want `log message must contain only English characters`

	// Только цифры и ASCII-пунктуация — ОК.
	slog.Info("200 OK, latency=15ms")

	// === SPECIAL CHARS ===

	// Каждый запрещённый символ по отдельности.
	slog.Info("hello@world")  // want `log message must not contain special characters or emoji`
	slog.Info("100% done")    // want `log message must not contain special characters or emoji`
	slog.Info("a^b")          // want `log message must not contain special characters or emoji`
	slog.Info("a&b")          // want `log message must not contain special characters or emoji`
	slog.Info("use * wisely") // want `log message must not contain special characters or emoji`
	slog.Info("path ~/home")  // want `log message must not contain special characters or emoji`
	slog.Info("foo#bar")      // want `log message must not contain special characters or emoji`
	slog.Info("cost $100")    // want `log message must not contain special characters or emoji`

	// Разрешённые символы: двоеточие, дефис, запятая, точка, скобки, слэш.
	slog.Info("host: 127.0.0.1, port: 8080")
	slog.Info("status - ok")
	slog.Info("path /api/v1/users")
	slog.Info("range (0, 100)")
	slog.Info("key=value")

	// Две точки — не эллипсис, ОК.
	slog.Info("range 1..10")

	// Ровно три точки — эллипсис, ловим.
	slog.Info("loading...") // want `log message must not contain special characters or emoji`

	// === SENSITIVE DATA ===

	x := "val"

	// Ключевое слово как подстрока — ловим.
	slog.Info("userPassword=" + x) // want `log message may contain sensitive data`

	// Регистр не важен.
	slog.Info("SECRET_VALUE: " + x) // want `log message may contain sensitive data`

	// Близкое слово, но не ключевое — НЕ ловим.
	slog.Info("passport: " + x)

	// Чистый литерал с ключевым словом — НЕ ловим
	// (это описательное сообщение, а не вывод значения).
	slog.Info("password has been changed")
	slog.Info("secret rotation completed")

	// === МНОЖЕСТВЕННЫЕ НАРУШЕНИЯ НА ОДНОЙ СТРОКЕ ===

	// Заглавная буква + спецсимвол — обе диагностики.
	slog.Info("Server started!") // want `log message must start with a lowercase letter` `log message must not contain special characters or emoji`

	// Кириллица + спецсимвол.
	slog.Info("ошибка!") // want `log message must contain only English characters` `log message must not contain special characters or emoji`
}
