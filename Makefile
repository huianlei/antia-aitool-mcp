.PHONY: build test clean run lint fmt help install-deps

# Binary name
BINARY_NAME=antia-aitool-mcp

# Build variables
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install-deps: ## Install Go dependencies
	go mod download
	go mod tidy

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/server

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/server

run: ## Run the server with example config
	go run ./cmd/server --config configs/config.yaml

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	go test -v -tags=integration ./tests/integration/...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	go clean

docker-build: ## Build Docker image
	docker build -t antia-aitool-mcp:$(VERSION) .

.DEFAULT_GOAL := help
