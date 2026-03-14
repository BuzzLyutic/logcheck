package analyzer

import "testing"

func TestSpecialCharsRule(t *testing.T) {
	rule := &SpecialCharsRule{}

	tests := []struct {
		name string
		lc   *LogCall
		want string
	}{
		{
			name: "чистое сообщение",
			lc:   &LogCall{MsgParts: []string{"server started"}},
			want: "",
		},
		{
			name: "восклицательный знак",
			lc:   &LogCall{MsgParts: []string{"server started!"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "тройной восклицательный",
			lc:   &LogCall{MsgParts: []string{"failed!!!"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "три точки (эллипсис)",
			lc:   &LogCall{MsgParts: []string{"loading..."}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "юникодный эллипсис (…)",
			lc:   &LogCall{MsgParts: []string{"loading\u2026"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "собака (@)",
			lc:   &LogCall{MsgParts: []string{"user@host"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "решётка (#)",
			lc:   &LogCall{MsgParts: []string{"issue#123"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "доллар ($)",
			lc:   &LogCall{MsgParts: []string{"cost $10"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "процент (%)",
			lc:   &LogCall{MsgParts: []string{"100% done"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "крышка (^)",
			lc:   &LogCall{MsgParts: []string{"a^b"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "амперсанд (&)",
			lc:   &LogCall{MsgParts: []string{"a&b"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "звёздочка (*)",
			lc:   &LogCall{MsgParts: []string{"use * wisely"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "тильда (~)",
			lc:   &LogCall{MsgParts: []string{"path ~/home"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "эмодзи ракета",
			lc:   &LogCall{MsgParts: []string{"deploying \U0001F680"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "эмодзи улыбка",
			lc:   &LogCall{MsgParts: []string{"all good \U0001F600"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "двоеточие и дефис — разрешены",
			lc:   &LogCall{MsgParts: []string{"status: ready - go"}},
			want: "",
		},
		{
			name: "одна точка — разрешена",
			lc:   &LogCall{MsgParts: []string{"host.example.com"}},
			want: "",
		},
		{
			name: "две точки — разрешены",
			lc:   &LogCall{MsgParts: []string{"range 1..10"}},
			want: "",
		},
		{
			name: "слэш — разрешён",
			lc:   &LogCall{MsgParts: []string{"path /api/v1"}},
			want: "",
		},
		{
			name: "запятая — разрешена",
			lc:   &LogCall{MsgParts: []string{"a, b, c"}},
			want: "",
		},
		{
			name: "знак равенства — разрешён",
			lc:   &LogCall{MsgParts: []string{"key=value"}},
			want: "",
		},
		{
			name: "скобки — разрешены",
			lc:   &LogCall{MsgParts: []string{"range (0, 100)"}},
			want: "",
		},
		{
			name: "пустые части",
			lc:   &LogCall{MsgParts: nil},
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

func TestIsEmoji(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"ракета", '\U0001F680', true},
		{"улыбка", '\U0001F600', true},
		{"солнце", '\u2600', true},
		{"ножницы (dingbat)", '\u2702', true},
		{"латинская A", 'A', false},
		{"цифра 1", '1', false},
		{"пробел", ' ', false},
		{"точка", '.', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmoji(tt.r); got != tt.want {
				t.Errorf("isEmoji(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}
