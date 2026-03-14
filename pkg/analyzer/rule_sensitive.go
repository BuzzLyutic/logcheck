package analyzer

import (
	"fmt"
	"strings"
)

// defaultSensitiveKeywords это список ключевых слов, которые указывают на то,
// что в журнал заносятся конфиденциальные данные.
var defaultSensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"secret",
	"token",
	"api_key", "apikey",
	"private_key", "privatekey",
	"access_key", "accesskey",
	"credential", "credentials",
}

// SensitiveRule проверяет, что сообщения журнала не содержат конфиденциальных данных.
//
// Он помечает только сообщения, которые являются конкатенациями (строковый литерал +
// переменная), поскольку этот шаблон предполагает, что значение регистрируется.
// Чистые строковые литералы, такие как "token validated", не помечаются.
type SensitiveRule struct {
	// Keywords переопределяет список ключевых слов по умолчанию, если он не пуст.
	Keywords []string
}

func (r *SensitiveRule) Name() string { return "sensitive" }

func (r *SensitiveRule) Check(lc *LogCall) string {
	// Чистый строковый литерал - разработчик ввел сообщение полностью, переменная не интерполируется.
	// Пропустить, чтобы избежать ложных срабатываний.
	if lc.MsgLit != "" {
		return ""
	}

	// Вообще никаких частей строки (чистый переменный аргумент) — анализировать нечего.
	if len(lc.MsgParts) == 0 {
		return ""
	}

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

	return ""
}

func (r *SensitiveRule) keywords() []string {
	if len(r.Keywords) > 0 {
		return r.Keywords
	}

	return defaultSensitiveKeywords
}
