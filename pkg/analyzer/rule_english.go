package analyzer

import "unicode"

// EnglishRule проверяет, что сообщение лога содержит только английские (ASCII) буквы.
// Non-ASCII буквы помечаются.
// Числа, ASCII пунктуация и пробелы допускаются.
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
