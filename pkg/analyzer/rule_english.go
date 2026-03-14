package analyzer

import "unicode"

// EnglishRule проверяет, что лог-сообщение содержит только английские (ASCII) символы.
// Не-ASCII буквы (кириллица, CJK, арабские и т.д.) считаются нарушением.
// Цифры, ASCII-пунктуация и пробелы разрешены.
type EnglishRule struct{}

func (r *EnglishRule) Name() string { return "english" }

func (r *EnglishRule) Check(lc *LogCall) string {
	for _, part := range lc.MsgParts {
		for _, ch := range part {
			if unicode.IsLetter(ch) && ch > unicode.MaxASCII {
				return "log message must contain only English characters"
			}
		}
	}

	return ""
}
