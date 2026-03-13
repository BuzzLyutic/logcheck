// Команда logcheck запускает линтер лог-сообщений.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"github.com/BuzzLyutic/logcheck/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
