// Пакет main реализует плагин для golangci-lint.
//
// Сборка:
//
//	CGO_ENABLED=1 go build -buildmode=plugin -o logcheck.so ./plugin/
//
// Использование в .golangci.yml:
//
//	linters-settings:
//	  custom:
//	    logcheck:
//	      path: ./logcheck.so
//	      description: "Линтер для проверки лог-сообщений"
package main

import (
	"github.com/BuzzLyutic/logcheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// analyzerPlugin реализует интерфейс golangci-lint для загрузки
// пользовательских анализаторов.
type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}

// AnalyzerPlugin — экспортируемая переменная, которую golangci-lint
// ищет при загрузке .so файла.
var AnalyzerPlugin analyzerPlugin //nolint:deadcode,unused
