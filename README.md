# logcheck

Линтер для Go, анализирующий лог-записи в коде и проверяющий их соответствие установленным правилам. Совместим с golangci-lint.

## Правила

| Правило | Описание |
|-|-|
| `lowercase` | Лог-сообщения должны начинаться со строчной буквы  |
| `english` | Лог-сообщения должны быть только на английском языке |
| `special-chars` |  Лог-сообщения не должны содержать спецсимволы или эмодзи |
| `sensitive` | Лог-сообщения не должны содержать потенциально чувствительные данные |

## Поддерживаемые логгеры

- `log/slog`
- `go.uber.org/zap`

## Установка

```bash
go install github.com/BuzzLyutic/logcheck/cmd/logcheck@latest
```

## Использование
```bash

# Standalone
logcheck ./...

# Via go vet
go vet -vettool=$(which logcheck) ./...
Build from source
Bash

git clone https://github.com/BuzzLyutic/logcheck.git
cd logcheck
make build
make test
```
