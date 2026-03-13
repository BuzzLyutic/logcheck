.PHONY: build test lint clean

BINARY := bin/logcheck

build:
	go build -o $(BINARY) ./cmd/logcheck

test:
	go test -v -count=1 ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out coverage.html

run: build
	./$(BINARY) ./...
