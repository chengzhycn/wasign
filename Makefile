# Makefile for go-cli-template
# Variables
BINARY_NAME=go-cli-template
MAIN_FILE=main.go
BUILD_DIR=bin
DOCKER_IMAGE=go-cli-template
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME} -w -s"

# Docker registry configuration (can be overridden by environment variables)
REGISTRY ?= ghcr.io
IMAGE_NAME ?= go-cli-template
FULL_IMAGE_NAME = $(if $(REGISTRY),$(REGISTRY)/$(IMAGE_NAME),$(IMAGE_NAME))

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker related variables
DOCKER_CMD=docker
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_RUN=$(DOCKER_CMD) run
DOCKER_TAG=$(DOCKER_CMD) tag
DOCKER_PUSH=$(DOCKER_CMD) push

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: ## Build the application for current platform
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

.PHONY: build-all
build-all: ## Build the application for multiple platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_FILE)

.PHONY: build-mac
build-mac: ## Build for macOS
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac $(MAIN_FILE)

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_FILE)

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Test targets
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race ./...

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

.PHONY: deps-vendor
deps-vendor: ## Vendor dependencies
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

# Linting and formatting
.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: vet
vet: ## Vet code
	@echo "Vetting code..."
	$(GOCMD) vet ./...

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	$(DOCKER_BUILD) \
		-t $(FULL_IMAGE_NAME):$(VERSION) .
	$(DOCKER_TAG) $(FULL_IMAGE_NAME):$(VERSION) $(FULL_IMAGE_NAME):latest

.PHONY: docker-build-multi
docker-build-multi: ## Build multi-platform Docker image
	@echo "Building multi-platform Docker image..."
	$(DOCKER_CMD) buildx create --use --name multi-platform-builder || true
	$(DOCKER_CMD) buildx build --platform linux/amd64,linux/arm64 \
		-t $(FULL_IMAGE_NAME):$(VERSION) \
		-t $(FULL_IMAGE_NAME):latest \
		--push .

.PHONY: docker-build-multi-local
docker-build-multi-local: ## Build multi-platform Docker image (local only)
	@echo "Building multi-platform Docker image (local only)..."
	$(DOCKER_CMD) buildx create --use --name multi-platform-builder || true
	$(DOCKER_CMD) buildx build --platform linux/amd64,linux/arm64 \
		-t $(FULL_IMAGE_NAME):$(VERSION) \
		-t $(FULL_IMAGE_NAME):latest \
		--load .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	$(DOCKER_RUN) --rm $(FULL_IMAGE_NAME):latest

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	@echo "Pushing Docker image..."
	$(DOCKER_PUSH) $(FULL_IMAGE_NAME):$(VERSION)
	$(DOCKER_PUSH) $(FULL_IMAGE_NAME):latest

.PHONY: docker-setup-buildx
docker-setup-buildx: ## Setup Docker Buildx for multi-platform builds
	@echo "Setting up Docker Buildx..."
	$(DOCKER_CMD) buildx create --name multi-platform-builder --use || true
	$(DOCKER_CMD) buildx inspect --bootstrap

# Release targets
.PHONY: release
release: clean build-all ## Create release builds
	@echo "Creating release builds..."
	@mkdir -p release
	@cd $(BUILD_DIR) && for file in *; do \
		tar -czf ../release/$$file.tar.gz $$file; \
	done
	@echo "Release artifacts created in release/ directory"

.PHONY: install
install: build ## Install the application
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

.PHONY: uninstall
uninstall: ## Uninstall the application
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)

# Development targets
.PHONY: dev
dev: ## Run in development mode
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_FILE)

.PHONY: dev-watch
dev-watch: ## Run with file watching (requires air)
	@echo "Running with file watching..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not found. Install with: go install github.com/cosmtrek/air@latest"; \
	fi

# Security targets
.PHONY: security-check
security-check: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: vuln-check
vuln-check: ## Check for vulnerabilities
	@echo "Checking for vulnerabilities..."
	$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Documentation targets
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Utility targets
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit Hash: $(COMMIT_HASH)"
	@echo "Build Time: $(BUILD_TIME)"

.PHONY: docker-config
docker-config: ## Show Docker configuration
	@echo "Docker Configuration:"
	@echo "  Registry: $(REGISTRY)"
	@echo "  Image Name: $(IMAGE_NAME)"
	@echo "  Full Image Name: $(FULL_IMAGE_NAME)"
	@echo "  Version: $(VERSION)"

.PHONY: check
check: fmt vet test lint ## Run all checks
	@echo "All checks completed"

.PHONY: ci
ci: deps check build-all ## Run CI pipeline
	@echo "CI pipeline completed"

# Install development tools
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/godoc@latest
	go install github.com/cosmtrek/air@latest 