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
			name: "только английский",
			lc:   &LogCall{MsgParts: []string{"starting server on port 8080"}},
			want: "",
		},
		{
			name: "кириллица",
			lc:   &LogCall{MsgParts: []string{"запуск сервера"}},
			want: "log message must contain only English characters",
		},
		{
			name: "смесь английского и кириллицы",
			lc:   &LogCall{MsgParts: []string{"server ошибка"}},
			want: "log message must contain only English characters",
		},
		{
			name: "китайские иероглифы",
			lc:   &LogCall{MsgParts: []string{"服务器启动"}},
			want: "log message must contain only English characters",
		},
		{
			name: "японская катакана",
			lc:   &LogCall{MsgParts: []string{"サーバー"}},
			want: "log message must contain only English characters",
		},
		{
			name: "акцентированная латиница (café)",
			lc:   &LogCall{MsgParts: []string{"café connection"}},
			want: "log message must contain only English characters",
		},
		{
			name: "немецкий умлаут (straße)",
			lc:   &LogCall{MsgParts: []string{"straße not found"}},
			want: "log message must contain only English characters",
		},
		{
			name: "только цифры и пунктуация",
			lc:   &LogCall{MsgParts: []string{"200 OK, latency=15ms"}},
			want: "",
		},
		{
			name: "пустые части",
			lc:   &LogCall{MsgParts: nil},
			want: "",
		},
		{
			name: "кириллица в конкатенации",
			lc: &LogCall{MsgParts: []string{
				"status: ",
				"ошибка",
			}},
			want: "log message must contain only English characters",
		},
		{
			name: "пустая строка в частях",
			lc:   &LogCall{MsgParts: []string{""}},
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
