package analyzer

import (
	"strings"
	"testing"
)

func TestSensitiveRule(t *testing.T) {
	rule := &SensitiveRule{}

	tests := []struct {
		name    string
		lc      *LogCall
		wantHas string // подстрока в диагностике ("" = нарушения нет)
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
			name:    "ключевое слово как подстрока (userPassword)",
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
		{
			name:    "только переменная (без частей)",
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
			name:    "дефолтное password НЕ срабатывает при кастомных",
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
