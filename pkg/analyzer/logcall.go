package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

// LogCall содержит информацию об обнаруженной функции логирования или вызове метода.
type LogCall struct {
	// Pos - это позиция выражения call.
	Pos token.Pos

	// Logger идентифицирует пакет логирования ("log/slog", "go.uber.org/zap").
	Logger string

	// Method - наименование метода логирования (e.g. "Info", "Error").
	Method string

	// MsgPos - позиция аргумента сообщения.
	MsgPos token.Pos

	// MsgLit - это строка сообщения без кавычек, если аргумент является одиночным
	// строковым литералом. В противном случае пустое значение.
	MsgLit string

	// MsgParts содержит каждый фрагмент строкового литерала, найденный в аргументе
	// сообщения, включая отдельные части выражений конкатенации.
	// Полезно для обнаружения конфиденциальных данных в выражениях типа "пароль: " + p.
	MsgParts []string
}

// extractStringLit возвращает значение строкового литерального выражения без кавычек.
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

// collectStringParts рекурсивно собирает все строковые литеральные фрагменты из
// expr, проходя через бинарные конкатенации "+" и заключенные в круглые скобки
// вложенные выражения.
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
