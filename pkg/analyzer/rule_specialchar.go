package analyzer

import "strings"

// forbiddenChars - это множество запрещенных символов в лог-сообщениях.
const forbiddenChars = "!@#$%^&*~"

// SpecialCharsRule проверяет, что лог-сообщение не содержит спец-символов или эмоджи.
type SpecialCharsRule struct{}

func (r *SpecialCharsRule) Name() string { return "special-chars" }

func (r *SpecialCharsRule) Check(lc *LogCall) string {
	for _, part := range lc.MsgParts {
		// Проверить наличие многоточий
		if strings.Contains(part, "...") || strings.ContainsRune(part, '…') {
			return "log message must not contain special characters or emoji"
		}

		for _, ch := range part {
			if strings.ContainsRune(forbiddenChars, ch) {
				return "log message must not contain special characters or emoji"
			}

			if isEmoji(ch) {
				return "log message must not contain special characters or emoji"
			}
		}
	}

	return ""
}

// isEmoji возвращает true если символ является эмоджи.
func isEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F: // Смайлики
		return true
	case r >= 0x1F300 && r <= 0x1F5FF: // Различные символы и пиктограммы
		return true
	case r >= 0x1F680 && r <= 0x1F6FF: // Транспортные и картографические обозначения
		return true
	case r >= 0x1F1E0 && r <= 0x1F1FF: // Флаги
		return true
	case r >= 0x2600 && r <= 0x26FF: // Разные символы
		return true
	case r >= 0x2700 && r <= 0x27BF: // Дингбаты
		return true
	case r >= 0x1F900 && r <= 0x1F9FF: // Дополнительные символы
		return true
	case r >= 0x1FA00 && r <= 0x1FAFF: // Символы Extended-A
		return true
	}

	return false
}
