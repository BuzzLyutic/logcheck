// Package analyzer реализует линтер для проверки лог-сообщений в Go.
//
// Анализирует вызовы поддерживаемых логгеров (log/slog, go.uber.org/zap)
// и проверяет сообщения на соответствие набору настраиваемых правил.
package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// configPath хранит путь к файлу конфигурации.
// Заполняется из флага -config при запуске.
var configPath string

// Analyzer — главная точка входа для линтера logcheck.
var Analyzer = &analysis.Analyzer{
	Name: "logcheck",
	Doc: "checks log messages for style and security issues: " +
		"lowercase start, English only, no special characters, no sensitive data",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.StringVar(
		&configPath,
		"config",
		"",
		"путь к файлу конфигурации (.logcheck.json)",
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Загружаем конфигурацию: из файла (если указан) или дефолтную.
	cfg := defaultConfig()
	if configPath != "" {
		var err error
		cfg, err = loadConfigFromFile(configPath)
		if err != nil {
			return nil, err
		}
	}

	// Собираем набор активных правил.
	rules, err := buildRules(cfg)
	if err != nil {
		return nil, err
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Нас интересуют только вызовы функций.
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		lc, ok := detectLogCall(pass, call)
		if !ok {
			return
		}

		// Прогоняем каждое правило.
		for _, r := range rules {
			msg := r.Check(lc)
			if msg == "" {
				continue
			}

			d := analysis.Diagnostic{
				Pos:     lc.MsgPos,
				Message: msg,
			}

			// Если правило умеет предлагать исправление — добавляем.
			if fr, ok := r.(FixableRule); ok {
				d.SuggestedFixes = fr.Fix(lc)
			}

			pass.Report(d)
		}
	})

	return nil, nil
}
