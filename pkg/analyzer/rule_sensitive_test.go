package analyzer

import (
	"regexp"
	"strings"
	"testing"
)

func TestSensitiveRule(t *testing.T) {
	// Дефолтное правило: стандартные ключевые слова, без паттернов.
	rule := &SensitiveRule{}

	tests := []struct {
		name    string
		lc      *LogCall
		wantHas string
	}{
		{
			name:    "чистый литерал с ключевым словом — ОК",
			lc:      &LogCall{MsgLit: "token validated", MsgParts: []string{"token validated"}},
			wantHas: "",
		},
		{
			name:    "литерал password reset — ОК",
			lc:      &LogCall{MsgLit: "password has been reset", MsgParts: []string{"password has been reset"}},
			wantHas: "",
		},
		{
			name:    "конкатенация с password",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"user password: "}},
			wantHas: `"password"`,
		},
		{
			name:    "конкатенация с token",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"token: "}},
			wantHas: `"token"`,
		},
		{
			name:    "конкатенация с api_key",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"api_key="}},
			wantHas: `"api_key"`,
		},
		{
			name:    "конкатенация с secret",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"secret: "}},
			wantHas: `"secret"`,
		},
		{
			name:    "конкатенация с credentials",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"credentials: "}},
			wantHas: `"credential"`,
		},
		{
			name:    "конкатенация с private_key",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"private_key="}},
			wantHas: `"private_key"`,
		},
		{
			name:    "ключевое слово в UPPER CASE",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"PASSWORD: "}},
			wantHas: `"password"`,
		},
		{
			name:    "подстрока userPassword",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"userPassword="}},
			wantHas: `"password"`,
		},
		{
			name:    "конкатенация без ключевых слов — ОК",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"user: "}},
			wantHas: "",
		},
		{
			name:    "похожее слово passport — ОК",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"passport: "}},
			wantHas: "",
		},
		{
			name:    "пустые части",
			lc:      &LogCall{MsgLit: "", MsgParts: nil},
			wantHas: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Check(tt.lc)
			if tt.wantHas == "" {
				if got != "" {
					t.Errorf("Check() = %q, ожидалось отсутствие диагностики", got)
				}
			} else {
				if !strings.Contains(got, tt.wantHas) {
					t.Errorf("Check() = %q, ожидалась подстрока %q", got, tt.wantHas)
				}
			}
		})
	}
}

func TestSensitiveRuleCustomKeywords(t *testing.T) {
	// Кастомные ключевые слова заменяют дефолтные.
	rule := &SensitiveRule{
		Keywords: []string{"ssn", "credit_card"},
	}

	tests := []struct {
		name    string
		lc      *LogCall
		wantHas string
	}{
		{
			name:    "кастомное ключевое слово ssn",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"user ssn: "}},
			wantHas: `"ssn"`,
		},
		{
			name:    "кастомное ключевое слово credit_card",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"credit_card="}},
			wantHas: `"credit_card"`,
		},
		{
			name:    "дефолтное password НЕ срабатывает",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"password: "}},
			wantHas: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Check(tt.lc)
			if tt.wantHas == "" {
				if got != "" {
					t.Errorf("Check() = %q, ожидалось отсутствие диагностики", got)
				}
			} else {
				if !strings.Contains(got, tt.wantHas) {
					t.Errorf("Check() = %q, ожидалась подстрока %q", got, tt.wantHas)
				}
			}
		})
	}
}

func TestSensitiveRulePatterns(t *testing.T) {
	// Только паттерны, ключевые слова отключены (пустой слайс).
	rule := &SensitiveRule{
		Keywords: []string{}, // Явно пустой — дефолтные не используются.
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)bearer\s+`),
			regexp.MustCompile(`(?i)auth[_-]?token`),
		},
	}

	tests := []struct {
		name    string
		lc      *LogCall
		wantHas string
	}{
		{
			name:    "паттерн bearer срабатывает",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"Authorization: Bearer "}},
			wantHas: "pattern",
		},
		{
			name:    "паттерн auth_token срабатывает",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"auth_token="}},
			wantHas: "pattern",
		},
		{
			name:    "паттерн authToken (без разделителя) срабатывает",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"authToken: "}},
			wantHas: "pattern",
		},
		{
			name:    "нет совпадения с паттерном — ОК",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"request id: "}},
			wantHas: "",
		},
		{
			name:    "чистый литерал — паттерн не проверяется",
			lc:      &LogCall{MsgLit: "bearer token received", MsgParts: []string{"bearer token received"}},
			wantHas: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Check(tt.lc)
			if tt.wantHas == "" {
				if got != "" {
					t.Errorf("Check() = %q, ожидалось отсутствие диагностики", got)
				}
			} else {
				if !strings.Contains(got, tt.wantHas) {
					t.Errorf("Check() = %q, ожидалась подстрока %q", got, tt.wantHas)
				}
			}
		})
	}
}
