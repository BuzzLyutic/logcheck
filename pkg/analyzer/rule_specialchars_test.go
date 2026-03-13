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
			name: "clean message",
			lc:   &LogCall{MsgParts: []string{"server started"}},
			want: "",
		},
		{
			name: "exclamation mark",
			lc:   &LogCall{MsgParts: []string{"server started!"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "multiple exclamation marks",
			lc:   &LogCall{MsgParts: []string{"connection failed!!!"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "ellipsis three dots",
			lc:   &LogCall{MsgParts: []string{"loading..."}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "unicode ellipsis",
			lc:   &LogCall{MsgParts: []string{"loading\u2026"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "at sign",
			lc:   &LogCall{MsgParts: []string{"user@host"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "emoji rocket",
			lc:   &LogCall{MsgParts: []string{"deploying \U0001F680"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "emoji smile",
			lc:   &LogCall{MsgParts: []string{"all good \U0001F600"}},
			want: "log message must not contain special characters or emoji",
		},
		{
			name: "colon and hyphen are allowed",
			lc:   &LogCall{MsgParts: []string{"status: ready - go"}},
			want: "",
		},
		{
			name: "single dot is allowed",
			lc:   &LogCall{MsgParts: []string{"connecting to host.example.com"}},
			want: "",
		},
		{
			name: "two dots are allowed",
			lc:   &LogCall{MsgParts: []string{"range 1..10"}},
			want: "",
		},
		{
			name: "empty parts",
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
		{"rocket", '\U0001F680', true},
		{"smile", '\U0001F600', true},
		{"sun", '\u2600', true},
		{"latin A", 'A', false},
		{"digit 1", '1', false},
		{"space", ' ', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmoji(tt.r); got != tt.want {
				t.Errorf("isEmoji(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}
