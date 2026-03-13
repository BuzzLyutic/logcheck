package analyzer

import "testing"

func TestEnglishRule(t *testing.T) {
	rule := &EnglishRule{}

	tests := []struct {
		name string
		lc   *LogCall
		want string
	}{
		{
			name: "english only",
			lc:   &LogCall{MsgParts: []string{"starting server on port 8080"}},
			want: "",
		},
		{
			name: "cyrillic",
			lc:   &LogCall{MsgParts: []string{"запуск сервера"}},
			want: "log message must contain only English characters",
		},
		{
			name: "mixed english and cyrillic",
			lc:   &LogCall{MsgParts: []string{"server ошибка"}},
			want: "log message must contain only English characters",
		},
		{
			name: "chinese characters",
			lc:   &LogCall{MsgParts: []string{"服务器启动"}},
			want: "log message must contain only English characters",
		},
		{
			name: "numbers and punctuation only",
			lc:   &LogCall{MsgParts: []string{"123-456: ok"}},
			want: "",
		},
		{
			name: "empty parts",
			lc:   &LogCall{MsgParts: nil},
			want: "",
		},
		{
			name: "cyrillic in concatenation part",
			lc: &LogCall{MsgParts: []string{
				"status: ",
				"ошибка",
			}},
			want: "log message must contain only English characters",
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
