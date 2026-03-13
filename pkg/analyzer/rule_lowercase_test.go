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
			name: "uppercase ASCII",
			lc:   &LogCall{MsgLit: "Starting server"},
			want: "log message must start with a lowercase letter",
		},
		{
			name: "lowercase ASCII",
			lc:   &LogCall{MsgLit: "starting server"},
			want: "",
		},
		{
			name: "starts with digit",
			lc:   &LogCall{MsgLit: "3 retries remaining"},
			want: "",
		},
		{
			name: "empty message",
			lc:   &LogCall{MsgLit: ""},
			want: "",
		},
		{
			name: "single uppercase letter",
			lc:   &LogCall{MsgLit: "A"},
			want: "log message must start with a lowercase letter",
		},
		{
			name: "concatenation without literal",
			lc:   &LogCall{MsgLit: "", MsgParts: []string{"hello"}},
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
