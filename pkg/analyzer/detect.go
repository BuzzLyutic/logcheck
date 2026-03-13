package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// logPackages перечисляет поддерживаемые пакеты логирования.
var logPackages = map[string]bool{
	"log/slog":        true,
	"go.uber.org/zap": true,
}

// logMethods отображает путь к пакету - имя метода - индекс аргумента сообщения.
var logMethods = map[string]map[string]int{
	"log/slog": {
		// Package-level functions and *Logger methods — msg is the first arg.
		"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
		// *Context variants — arg 0 is context.Context, msg is arg 1.
		"DebugContext": 1, "InfoContext": 1, "WarnContext": 1, "ErrorContext": 1,
		// Log — arg 0 is context, arg 1 is level, msg is arg 2.
		"Log": 2,
	},
	"go.uber.org/zap": {
		// *zap.Logger — msg is the first arg.
		"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
		"DPanic": 0, "Panic": 0, "Fatal": 0,
		// *zap.SugaredLogger — "w" (key-value) variants.
		"Debugw": 0, "Infow": 0, "Warnw": 0, "Errorw": 0,
		"DPanicw": 0, "Panicw": 0, "Fatalw": 0,
		// *zap.SugaredLogger — "f" (format) variants.
		"Debugf": 0, "Infof": 0, "Warnf": 0, "Errorf": 0,
		"DPanicf": 0, "Panicf": 0, "Fatalf": 0,
	},
}

// detectLogCall проверяет, является ли вызов известным лог-вызовом, и, если да,
// извлекает аргумент с сообщением в LogCall.
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

	// Попытка получить полное буквальное значение (простой строковый литерал).
	lc.MsgLit, _ = extractStringLit(msgArg)

	// Сборка всех фрагментов строки.
	collectStringParts(msgArg, &lc.MsgParts)

	return lc, true
}

// resolvePackagePath возвращает путь импорта пакета, которому принадлежит
// метод/функция, на которые ссылается sel, но только если это поддерживаемый логгер.
func resolvePackagePath(pass *analysis.Pass, sel *ast.SelectorExpr) string {
	if obj := pass.TypesInfo.Uses[sel.Sel]; obj != nil {
		if fn, ok := obj.(*types.Func); ok {
			if pkg := fn.Pkg(); pkg != nil && logPackages[pkg.Path()] {
				return pkg.Path()
			}
		}
	}

	// Вызов метода для типизированного получателя: logger.Info, sugar.Warnw, etc.
	if selection, ok := pass.TypesInfo.Selections[sel]; ok {
		if fn, ok := selection.Obj().(*types.Func); ok {
			if pkg := fn.Pkg(); pkg != nil && logPackages[pkg.Path()] {
				return pkg.Path()
			}
		}
	}

	return ""
}
