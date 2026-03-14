package analyzer_test

import (
	"testing"

	"github.com/BuzzLyutic/logcheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

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

func TestAnalyzerEdgeCases(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "edgecases")
}

func TestAnalyzerSlogContext(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "slogcontext")
}

func TestAnalyzerMultiline(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "multiline")
}

func TestAnalyzerFixLowercase(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer.Analyzer, "fixlowercase")
}
