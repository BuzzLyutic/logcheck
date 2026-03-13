// Пакет analyzer реализует линтер для лог-сообщений.
package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer - это точка входа для logcheck линтера.
var Analyzer = &analysis.Analyzer{
	Name: "logcheck",
	Doc: "checks log messages for style and security issues: " +
		"lowercase start, English only, no special characters, no sensitive data",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// rules это множество активных проверок.
var rules []Rule

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		lc, ok := detectLogCall(pass, call)
		if !ok {
			return
		}

		for _, r := range rules {
			if msg := r.Check(lc); msg != "" {
				pass.Reportf(lc.MsgPos, "%s", msg)
			}
		}
	})

	return nil, nil
}