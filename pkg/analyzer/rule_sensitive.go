package analyzer

import (
	"fmt"
	"regexp"
	"strings"
)

// defaultSensitiveKeywords — ключевые слова, указывающие на
// потенциально чувствительные данные.
var defaultSensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"secret",
	"token",
	"api_key", "apikey",
	"private_key", "privatekey",
	"access_key", "accesskey",
	"credential", "credentials",
}

// SensitiveRule проверяет, что лог-сообщения не выводят чувствительные данные.
//
// Срабатывает только на конкатенациях (строковый литерал + переменная),
// потому что такой паттерн означает вывод значения переменной.
// Чистые литералы вроде "token validated" не считаются нарушением.
type SensitiveRule struct {
	// Keywords переопределяет список ключевых слов (если не пуст).
	Keywords []string

	// Patterns — скомпилированные регулярные выражения для поиска
	// чувствительных данных. Проверяются в дополнение к ключевым словам.
	Patterns []*regexp.Regexp
}

func (r *SensitiveRule) Name() string { return "sensitive" }

func (r *SensitiveRule) Check(lc *LogCall) string {
	// Чистый строковый литерал — разработчик написал описательное сообщение,
	// переменная не интерполируется. Пропускаем, чтобы избежать
	// ложных срабатываний вроде "token validated".
	if lc.MsgLit != "" {
		return ""
	}

	// Нет строковых частей (чистая переменная) — нечего анализировать.
	if len(lc.MsgParts) == 0 {
		return ""
	}

	// Проверка по ключевым словам.
	keywords := r.keywords()
	for _, part := range lc.MsgParts {
		lower := strings.ToLower(part)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				return fmt.Sprintf(
					"log message may contain sensitive data (keyword %q)", kw,
				)
			}
		}
	}

	// Проверка по регулярным выражениям (кастомные паттерны).
	for _, part := range lc.MsgParts {
		for _, pat := range r.Patterns {
			if pat.MatchString(part) {
				return fmt.Sprintf(
					"log message may contain sensitive data (pattern %q)", pat.String(),
				)
			}
		}
	}

	return ""
}

func (r *SensitiveRule) keywords() []string {
	// Если Keywords явно задан (даже пустой слайс) — используем его.
	// Это позволяет отключить ключевые слова, оставив только паттерны.
	if r.Keywords != nil {
		return r.Keywords
	}

	return defaultSensitiveKeywords
}
