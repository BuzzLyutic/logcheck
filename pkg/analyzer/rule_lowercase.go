package analyzer

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// LowercaseRule проверяет, что лог-сообщение начинается со строчной буквы.
type LowercaseRule struct{}

func (r *LowercaseRule) Name() string { return "lowercase" }

func (r *LowercaseRule) Check(lc *LogCall) string {
	if lc.MsgLit == "" {
		return ""
	}

	ch, size := utf8.DecodeRuneInString(lc.MsgLit)
	if size == 0 {
		return ""
	}

	if unicode.IsUpper(ch) {
		return "log message must start with a lowercase letter"
	}

	return ""
}

// Fix предлагает заменить первую заглавную букву на строчную.
// Работает только для ASCII-букв — для многобайтных символов (Ü, Ö)
// автоисправление не предлагается, потому что размер руны в исходнике
// может отличаться (escape-последовательности и т.д.).
func (r *LowercaseRule) Fix(lc *LogCall) []analysis.SuggestedFix {
	if lc.MsgLit == "" || lc.MsgRaw == "" {
		return nil
	}

	ch, _ := utf8.DecodeRuneInString(lc.MsgLit)
	if !unicode.IsUpper(ch) {
		return nil
	}

	// Автоисправление только для ASCII (1 байт в исходнике = 1 байт замены).
	if ch > unicode.MaxASCII {
		return nil
	}

	lower := unicode.ToLower(ch)

	return []analysis.SuggestedFix{
		{
			Message: "заменить первую букву на строчную",
			TextEdits: []analysis.TextEdit{
				{
					Pos:     lc.MsgPos + 1, // +1 — пропускаем открывающую кавычку
					End:     lc.MsgPos + 2, // ASCII = 1 байт
					NewText: []byte(string(lower)),
				},
			},
		},
	}
}
