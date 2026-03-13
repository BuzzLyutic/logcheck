package analyzer

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestExtractStringLit(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want string
		ok   bool
	}{
		{
			name: "double-quoted string",
			expr: &ast.BasicLit{Kind: token.STRING, Value: `"hello world"`},
			want: "hello world",
			ok:   true,
		},
		{
			name: "raw string",
			expr: &ast.BasicLit{Kind: token.STRING, Value: "`raw string`"},
			want: "raw string",
			ok:   true,
		},
		{
			name: "string with escapes",
			expr: &ast.BasicLit{Kind: token.STRING, Value: `"line\nbreak"`},
			want: "line\nbreak",
			ok:   true,
		},
		{
			name: "empty string",
			expr: &ast.BasicLit{Kind: token.STRING, Value: `""`},
			want: "",
			ok:   true,
		},
		{
			name: "integer literal",
			expr: &ast.BasicLit{Kind: token.INT, Value: "42"},
			want: "",
			ok:   false,
		},
		{
			name: "identifier",
			expr: &ast.Ident{Name: "variable"},
			want: "",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := extractStringLit(tt.expr)
			if ok != tt.ok {
				t.Errorf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCollectStringParts(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want []string
	}{
		{
			name: "single literal",
			expr: &ast.BasicLit{Kind: token.STRING, Value: `"hello"`},
			want: []string{"hello"},
		},
		{
			name: "two literals concatenated",
			expr: &ast.BinaryExpr{
				X:  &ast.BasicLit{Kind: token.STRING, Value: `"hello "`},
				Op: token.ADD,
				Y:  &ast.BasicLit{Kind: token.STRING, Value: `"world"`},
			},
			want: []string{"hello ", "world"},
		},
		{
			name: "literal plus variable",
			expr: &ast.BinaryExpr{
				X:  &ast.BasicLit{Kind: token.STRING, Value: `"user: "`},
				Op: token.ADD,
				Y:  &ast.Ident{Name: "name"},
			},
			want: []string{"user: "},
		},
		{
			name: "triple concatenation",
			expr: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Kind: token.STRING, Value: `"a"`},
					Op: token.ADD,
					Y:  &ast.BasicLit{Kind: token.STRING, Value: `"b"`},
				},
				Op: token.ADD,
				Y:  &ast.BasicLit{Kind: token.STRING, Value: `"c"`},
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "parenthesised literal",
			expr: &ast.ParenExpr{
				X: &ast.BasicLit{Kind: token.STRING, Value: `"wrapped"`},
			},
			want: []string{"wrapped"},
		},
		{
			name: "variable only",
			expr: &ast.Ident{Name: "msg"},
			want: nil,
		},
		{
			name: "non-string binary",
			expr: &ast.BinaryExpr{
				X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
				Op: token.ADD,
				Y:  &ast.BasicLit{Kind: token.INT, Value: "2"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			collectStringParts(tt.expr, &got)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d parts %v, want %d parts %v", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("part[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
