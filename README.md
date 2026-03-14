# logcheck

Go-линтер для проверки лог-сообщений. Совместим с golangci-lint.

Анализирует вызовы логгеров и проверяет сообщения на соответствие правилам
оформления и безопасности.

[![CI](https://github.com/BuzzLyutic/logcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/BuzzLyutic/logcheck/actions/workflows/ci.yml)

## Правила

| Правило | Описание | Пример нарушения |
|---------|----------|------------------|
| `lowercase` | Сообщение должно начинаться со строчной буквы | `slog.Info("Starting server")` |
| `english` | Сообщение должно быть на английском языке | `slog.Info("запуск сервера")` |
| `special-chars` | Нет спецсимволов (`!@#$%^&*~`), эллипсисов и эмодзи | `slog.Info("started!")` |
| `sensitive` | Нет чувствительных данных в конкатенациях | `slog.Info("token: " + t)` |

## Поддерживаемые логгеры

- `log/slog` — все методы (`Info`, `Error`, `Warn`, `Debug`, `InfoContext`, `Log`, ...)
- `go.uber.org/zap` — `*Logger` и `*SugaredLogger`

## Установка

### Standalone бинарник

```bash
go install github.com/BuzzLyutic/logcheck/cmd/logcheck@latest
```

### Из исходников
```bash

git clone https://github.com/BuzzLyutic/logcheck.git
cd logcheck
make build
```

## Использование

### Standalone

```bash

# Проверить все пакеты в проекте
logcheck ./...

# Проверить конкретный пакет
logcheck ./internal/server/
```

### Через go vet

```bash
go vet -vettool=$(which logcheck) ./...
```

### Как плагин golangci-lint

Собрать плагин:

```bash

cd logcheck
CGO_ENABLED=1 go build -buildmode=plugin -o logcheck.so ./plugin/
```
Добавить в .golangci.yml проекта:

```YAML

linters-settings:
  custom:
    logcheck:
      path: /path/to/logcheck.so
      description: "Log message linter"
      original-url: github.com/BuzzLyutic/logcheck
```
Запустить:
```bash
golangci-lint run ./...
```

### Пример вывода
```text
$ logcheck ./examples/demo/
examples/demo/main.go:10:12: log message must start with a lowercase letter
examples/demo/main.go:13:13: log message must contain only English characters
examples/demo/main.go:16:12: log message must not contain special characters or emoji
examples/demo/main.go:20:12: log message may contain sensitive data (keyword "token")
```

## Разработка
```bash

# Тесты
make test

# Покрытие
make cover

# Линтинг проекта
make lint

# Сборка
make build
```

## Архитектура

```text

cmd/logcheck/          — точка входа standalone-бинарника
pkg/analyzer/
  analyzer.go          — главный анализатор (обход AST, запуск правил)
  detect.go            — определение вызовов логгеров по типовой информации
  logcall.go           — модель лог-вызова, извлечение строковых аргументов
  rule.go              — интерфейс Rule
  rule_lowercase.go    — правило: строчная буква
  rule_english.go      — правило: только английский
  rule_specialchars.go — правило: спецсимволы и эмодзи
  rule_sensitive.go    — правило: чувствительные данные
  testdata/            — тестовые пакеты для analysistest
plugin/                — плагин для golangci-lint
examples/              — примеры для демонстрации
```

## Ограничения
Анализируются только строковые литералы. Если сообщение передаётся через переменную или константу, линтер его пропускает.
Правило sensitive срабатывает только на конкатенациях ("password: " + p). Чистые литералы вроде "password reset" не считаются нарушением.
