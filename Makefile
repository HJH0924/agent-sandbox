.PHONY: help build run test lint format proto clean

# Variables
BINARY_NAME=sandbox-server
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

run: ## Run the server
	$(GO) run cmd/server/main.go -c configs/config.yaml

test: ## Run tests
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run --timeout=5m

format: ## Format code
	$(GO) fmt ./...
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -w .

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf gen/
	rm -f coverage.out coverage.html
	$(GO) clean

generate: ## Generate protobuf code
	buf generate

docs: ## Run docs
	cd docs && pnpm run docs:dev

docker-build: ## Build Docker image
	docker build -t agent-sandbox:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 -v $(PWD)/configs:/app/configs agent-sandbox:latest

.DEFAULT_GOAL := help

