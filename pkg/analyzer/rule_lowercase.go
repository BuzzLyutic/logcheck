package analyzer

import (
	"unicode"
	"unicode/utf8"
)

// LowercaseRule проверяет, что сообщение лога начинается с маленькой буквы.
type LowercaseRule struct{}

func (r *LowercaseRule) Name() string { return "lowercase" }

func (r *LowercaseRule) Check(lc *LogCall) string {
	// Если сообщение — не простой строковый литерал (переменная,
	// конкатенация, начинающаяся с переменной), мы не можем
	// определить первый символ — пропускаем.
	if lc.MsgLit == "" {
		return ""
	}

	// Декодируем первый символ (руну) из строки.
	// utf8.DecodeRuneInString безопаснее и понятнее, чем цикл for-range
	// с return на первой итерации.
	ch, size := utf8.DecodeRuneInString(lc.MsgLit)
	if size == 0 {
		// Пустая строка — нечего проверять.
		return ""
	}

	if unicode.IsUpper(ch) {
		return "log message must start with a lowercase letter"
	}

	return ""
}
