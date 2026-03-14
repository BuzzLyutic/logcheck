package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// logPackages — множество поддерживаемых пакетов логирования.
var logPackages = map[string]bool{
	"log/slog":        true,
	"go.uber.org/zap": true,
}

// logMethods — для каждого пакета: имя метода - индекс аргумента-сообщения.
// Индекс нужен, потому что у некоторых методов перед сообщением идут
// другие аргументы (context, level и т.д.).
var logMethods = map[string]map[string]int{
	"log/slog": {
		"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
		"DebugContext": 1, "InfoContext": 1, "WarnContext": 1, "ErrorContext": 1,
		"Log": 2,
	},
	"go.uber.org/zap": {
		"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
		"DPanic": 0, "Panic": 0, "Fatal": 0,
		"Debugw": 0, "Infow": 0, "Warnw": 0, "Errorw": 0,
		"DPanicw": 0, "Panicw": 0, "Fatalw": 0,
		"Debugf": 0, "Infof": 0, "Warnf": 0, "Errorf": 0,
		"DPanicf": 0, "Panicf": 0, "Fatalf": 0,
	},
}

// detectLogCall проверяет, является ли call вызовом поддерживаемого логгера.
// Если да — извлекает информацию о вызове в LogCall.
func detectLogCall(pass *analysis.Pass, call *ast.CallExpr) (*LogCall, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	methodName := sel.Sel.Name

	pkgPath := resolvePackagePath(pass, sel)
	if pkgPath == "" {
		return nil, false
	}

	methods, ok := logMethods[pkgPath]
	if !ok {
		return nil, false
	}

	msgIdx, ok := methods[methodName]
	if !ok {
		return nil, false
	}

	if msgIdx >= len(call.Args) {
		return nil, false
	}

	msgArg := call.Args[msgIdx]

	lc := &LogCall{
		Pos:    call.Pos(),
		Logger: pkgPath,
		Method: methodName,
		MsgPos: msgArg.Pos(),
	}

	// Извлекаем значение, если это простой строковый литерал.
	lc.MsgLit, _ = extractStringLit(msgArg)

	// Сохраняем исходный текст литерала (с кавычками) для автоисправления.
	if lit, ok := msgArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		lc.MsgRaw = lit.Value
	}

	// Собираем все строковые фрагменты (включая части конкатенации).
	collectStringParts(msgArg, &lc.MsgParts)

	return lc, true
}

// resolvePackagePath определяет путь пакета, которому принадлежит
// метод/функция, вызываемая через sel.
// Возвращает пустую строку, если это не поддерживаемый логгер.
func resolvePackagePath(pass *analysis.Pass, sel *ast.SelectorExpr) string {
	// Случай 1: прямой вызов функции пакета (slog.Info).
	if obj := pass.TypesInfo.Uses[sel.Sel]; obj != nil {
		if fn, ok := obj.(*types.Func); ok {
			if pkg := fn.Pkg(); pkg != nil && logPackages[pkg.Path()] {
				return pkg.Path()
			}
		}
	}

	// Случай 2: вызов метода на переменной (logger.Info).
	if selection, ok := pass.TypesInfo.Selections[sel]; ok {
		if fn, ok := selection.Obj().(*types.Func); ok {
			if pkg := fn.Pkg(); pkg != nil && logPackages[pkg.Path()] {
				return pkg.Path()
			}
		}
	}

	return ""
}
