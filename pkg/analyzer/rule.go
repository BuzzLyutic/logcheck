package analyzer

import "golang.org/x/tools/go/analysis"

// Rule — интерфейс одного правила проверки, применяемого
// к каждому обнаруженному вызову логгера.
type Rule interface {
	// Name возвращает короткий уникальный идентификатор правила (например "lowercase").
	Name() string

	// Check проверяет lc и возвращает текст диагностики при нарушении.
	// Пустая строка означает, что нарушения нет.
	Check(lc *LogCall) string
}

// FixableRule — расширение Rule для правил, умеющих предлагать
// автоматическое исправление.
type FixableRule interface {
	Rule

	// Fix возвращает список предложенных исправлений для нарушения.
	// Вызывается только когда Check вернул непустую строку.
	Fix(lc *LogCall) []analysis.SuggestedFix
}
