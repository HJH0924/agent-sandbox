.PHONY: help install-deps install-hooks format lint build build-linux run test clean generate docs docker shell e2b

# Variables
BINARY_NAME=api-server
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-deps: ## Check and install development dependencies
	@command -v golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest)
	@command -v gofumpt > /dev/null || (echo "Installing gofumpt..." && go install mvdan.cc/gofumpt@latest)
	@command -v goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@command -v goimports-reviser > /dev/null || (echo "Installing goimports-reviser..." && go install github.com/incu6us/goimports-reviser/v3@latest)
	@command -v buf > /dev/null || (echo "Installing buf..." && go install github.com/bufbuild/buf/cmd/buf@latest)
	@$(MAKE) install-hooks
	@echo ""
	@echo "✅ All development dependencies installed successfully!"

install-hooks: ## Install git pre-commit hooks
	@if [ -d .git ]; then \
		echo "Installing pre-commit hook..."; \
		mkdir -p .git/hooks; \
		cp scripts/pre-commit .git/hooks/pre-commit; \
		chmod +x .git/hooks/pre-commit; \
		echo "✅ Pre-commit hook installed successfully!"; \
	else \
		echo "⚠️  Not a git repository, skipping hook installation"; \
	fi

format: ## Format code
	gofumpt -l -w .
	goimports-reviser -imports-order=std,project,company,general  -recursive ./

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest)
	golangci-lint run -c .golangci.toml ./...

build: format lint ## Build the binary
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) .

build-linux: format lint ## Build the binary for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) .

run: format lint ## Run the server
	$(GO) run cmd/server/main.go -c configs/config.yaml

test: format lint ## Run tests
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf gen/
	rm -f coverage.out coverage.html
	$(GO) clean

generate: ## Generate protobuf code
	buf generate

docs: ## Run docs
	cd docs && pnpm run docs:dev

docker: build-linux ## Build program and run Docker container locally
	docker-compose down
	docker-compose up -d --build
	@echo ""
	@echo "✅ Docker container started successfully!"

shell: ## Connect to running Docker container shell
	@docker-compose exec agent-sandbox /bin/bash

e2b: build-linux ## Build E2B sandbox template
	e2b template build \
		--dockerfile e2b.Dockerfile \
		--name agent-sandbox \
		--cpu-count 2 \
		--memory-mb 2048 \
		--cmd "/app/bin/api-server -c /app/configs/config.yaml"

.DEFAULT_GOAL := help

