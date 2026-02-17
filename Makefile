.PHONY: help build run test lint fmt clean install

help:
	@echo "asana-cli - Asana CLI with TUI and sync daemon"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  install     - Install the binary"
	@echo "  run         - Run the CLI"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format code"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download dependencies"

build:
	go build -o asana-cli main.go

install:
	go install

run:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

clean:
	rm -f asana-cli
	go clean

deps:
	go mod download
	go mod tidy

dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest