package analyzer

import "testing"

func TestLowercaseRule(t *testing.T) {
	rule := &LowercaseRule{}

	tests := []struct {
		name string
		lc   *LogCall
		want string
	}{
		{
			name: "заглавная ASCII",
			lc:   &LogCall{MsgLit: "Starting server"},
			want: "log message must start with a lowercase letter",
		},
		{
			name: "строчная ASCII",
			lc:   &LogCall{MsgLit: "starting server"},
			want: "",
		},
		{
			name: "начинается с цифры",
			lc:   &LogCall{MsgLit: "3 retries remaining"},
			want: "",
		},
		{
			name: "пустое сообщение",
			lc:   &LogCall{MsgLit: ""},
			want: "",
		},
		{
			name: "одна заглавная буква",
			lc:   &LogCall{MsgLit: "A"},
			want: "log message must start with a lowercase letter",
		},
		{
			name: "одна строчная буква",
			lc:   &LogCall{MsgLit: "a"},
			want: "",
		},
		{
			name: "начинается с пунктуации",
			lc:   &LogCall{MsgLit: "[server] ready"},
			want: "",
		},
		{
			name: "заглавная не-ASCII (Ü)",
			lc:   &LogCall{MsgLit: "Über alles"},
			want: "log message must start with a lowercase letter",
		},
		{
			name: "строчная не-ASCII (ü)",
			lc:   &LogCall{MsgLit: "über alles"},
			want: "",
		},
		{
			name: "конкатенация (MsgLit пуст) — пропускаем",
			lc:   &LogCall{MsgLit: "", MsgParts: []string{"Hello"}},
			want: "",
		},
		{
			name: "пробел в начале",
			lc:   &LogCall{MsgLit: " starting"},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Check(tt.lc)
			if got != tt.want {
				t.Errorf("Check() = %q, want %q", got, tt.want)
			}
		})
	}
}
