// Команда logcheck запускает линтер лог-сообщений.
package main

import (
	"github.com/BuzzLyutic/logcheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
