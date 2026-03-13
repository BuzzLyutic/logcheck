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
		wantHas string // подстрока, которая должна появиться в диагностике ("" = no diagnostic)
	}{
		{
			name:    "pure literal with keyword is OK",
			lc:      &LogCall{MsgLit: "token validated", MsgParts: []string{"token validated"}},
			wantHas: "",
		},
		{
			name:    "pure literal password reset is OK",
			lc:      &LogCall{MsgLit: "password reset email sent", MsgParts: []string{"password reset email sent"}},
			wantHas: "",
		},
		{
			name:    "concatenation with password keyword",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"user password: "}},
			wantHas: `"password"`,
		},
		{
			name:    "concatenation with token keyword",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"token: "}},
			wantHas: `"token"`,
		},
		{
			name:    "concatenation with api_key",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"api_key="}},
			wantHas: `"api_key"`,
		},
		{
			name:    "concatenation with secret",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"secret: "}},
			wantHas: `"secret"`,
		},
		{
			name:    "keyword is case-insensitive",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"PASSWORD: "}},
			wantHas: `"password"`,
		},
		{
			name:    "concatenation without sensitive keyword",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"user: "}},
			wantHas: "",
		},
		{
			name:    "empty parts",
			lc:      &LogCall{MsgLit: "", MsgParts: nil},
			wantHas: "",
		},
		{
			name:    "pure variable no parts",
			lc:      &LogCall{MsgLit: "", MsgParts: nil},
			wantHas: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Check(tt.lc)
			if tt.wantHas == "" {
				if got != "" {
					t.Errorf("Check() = %q, want no diagnostic", got)
				}
			} else {
				if !strings.Contains(got, tt.wantHas) {
					t.Errorf("Check() = %q, want it to contain %q", got, tt.wantHas)
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
			name:    "custom keyword ssn",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"user ssn: "}},
			wantHas: `"ssn"`,
		},
		{
			name:    "default keyword password not flagged",
			lc:      &LogCall{MsgLit: "", MsgParts: []string{"password: "}},
			wantHas: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Check(tt.lc)
			if tt.wantHas == "" {
				if got != "" {
					t.Errorf("Check() = %q, want no diagnostic", got)
				}
			} else {
				if !strings.Contains(got, tt.wantHas) {
					t.Errorf("Check() = %q, want it to contain %q", got, tt.wantHas)
				}
			}
		})
	}
}
