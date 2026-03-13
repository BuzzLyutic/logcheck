// Пакет analyzer реализует линтер для лог-сообщений.
package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

// Analyzer является точкой входа для logcheck линтера.
var Analyzer = newAnalyzer()

func newAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "logcheck",
		Doc:  "checks log messages for common issues",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	return nil, nil
}
