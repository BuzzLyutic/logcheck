// Этот файл проверяет edge-cases со строковыми литералами:
// raw-строки, конкатенация констант, переменные.
package multiline

import "log/slog"

const myMsg = "Starting background job"

func multiline() {
	// Raw-строка (backtick) с заглавной буквой.
	slog.Info(`Starting raw server`) // want `log message must start with a lowercase letter`

	// Raw-строка, всё ОК.
	slog.Info(`starting raw server`)

	// Константа — анализатор видит значение литерала в месте вызова?
	// Нет: константа подставляется компилятором, в AST остаётся Ident.
	// Линтер не может извлечь значение из Ident — пропускает.
	// Это осознанное ограничение.
	slog.Info(myMsg)

	// Переменная — тоже пропускаем, значение неизвестно.
	msg := getMessage()
	slog.Info(msg)

	// Конкатенация двух литералов — MsgLit будет пустым (это BinaryExpr),
	// но MsgParts соберёт оба фрагмента.
	// Lowercase правило проверяет только MsgLit — пропустит.
	// English и SpecialChars проверяют MsgParts — поймают.
	slog.Info("starting " + "сервер") // want `log message must contain only English characters`
}

func getMessage() string {
	return "some message"
}
