package analyzer_test

import (
	"testing"
	"golang.org/x/tools/go/analysis/analysistest"
	"github.com/BuzzLyutic/logcheck/pkg/analyzer"
)

// TestAnalyzerSmoke проверяет, что анализатор работает без ошибок в коде, содержащем вызовы slog.
func TestAnalyzerSmoke(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "slogbasic")
}

func TestAnalyzerSlogRules(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "slogrules")
}

func TestAnalyzerZapRules(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "zaprules")
}
