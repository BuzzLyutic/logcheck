package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

// LogCall содержит информацию об обнаруженном вызове функции логирования.
type LogCall struct {
	// Pos — позиция всего выражения вызова в исходном файле.
	Pos token.Pos

	// Logger — путь пакета логирования (например "log/slog", "go.uber.org/zap").
	Logger string

	// Method — имя метода логирования (например "Info", "Error").
	Method string

	// MsgPos — позиция аргумента-сообщения в исходном файле.
	// Для строкового литерала это позиция открывающей кавычки.
	MsgPos token.Pos

	// MsgLit — строковое значение сообщения (без кавычек),
	// если аргумент — простой строковый литерал.
	// Пустая строка, если аргумент — переменная или конкатенация.
	MsgLit string

	// MsgRaw — исходный текст строкового литерала вместе с кавычками
	// (например `"Starting server"` или `` `raw string` ``).
	// Нужен для автоисправления — чтобы знать точные позиции символов.
	// Пустой, если аргумент — не простой строковый литерал.
	MsgRaw string

	// MsgParts — все строковые фрагменты, найденные в аргументе-сообщении.
	// Включает отдельные части конкатенации ("hello " + var + " world"
	// даст ["hello ", " world"]).
	// Полезно для поиска чувствительных данных в составных выражениях.
	MsgParts []string
}

// extractStringLit извлекает строковое значение из литерала (без кавычек).
func extractStringLit(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return val, true
}

// collectStringParts рекурсивно собирает все строковые фрагменты
// из выражения, проходя через бинарные операции "+" и скобки.
func collectStringParts(expr ast.Expr, parts *[]string) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if val, err := strconv.Unquote(e.Value); err == nil {
				*parts = append(*parts, val)
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			collectStringParts(e.X, parts)
			collectStringParts(e.Y, parts)
		}
	case *ast.ParenExpr:
		collectStringParts(e.X, parts)
	}
}
