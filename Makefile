.PHONY: build clean install test run help

# Build variables
BINARY_NAME=atlas-migrate-status
BUILD_DIR=bin
GO=go

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install .
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test: ## Run tests
	$(GO) test -v ./...

run: build ## Build and run with example (requires DB_URL env var)
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL environment variable not set"; \
		echo "Example: make run DB_URL='postgres://user:pass@localhost:5432/dbname'"; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME) --url "$(DB_URL)"

fmt: ## Format code
	$(GO) fmt ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

.DEFAULT_GOAL := help
