package analyzer_test

import (
	"testing"
	"golang.org/x/tools/go/analysis/analysistest"
	"github.com/BuzzLyutic/logcheck/pkg/analyzer"
)

func TestAnalyzerSkeleton(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "example")
}
