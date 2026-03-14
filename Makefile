.PHONY: build build-plugin test cover lint clean run

BINARY := bin/logcheck
PLUGIN := bin/logcheck.so

# Сборка standalone-бинарника.
build:
	go build -o $(BINARY) ./cmd/logcheck

# Сборка плагина для golangci-lint.
# Требует CGO_ENABLED=1 и Linux/macOS.
build-plugin:
	CGO_ENABLED=1 go build -buildmode=plugin -o $(PLUGIN) ./plugin/

# Запуск всех тестов.
test:
	go test -v -count=1 ./...

# Покрытие кода тестами.
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт: coverage.html"

# Запуск линтеров на самом проекте.
lint:
	golangci-lint run ./...

# Удаление артефактов сборки.
clean:
	rm -rf bin/ coverage.out coverage.html

# Сборка и запуск на текущем проекте.
run: build
	./$(BINARY) ./...
