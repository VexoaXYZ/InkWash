.PHONY: build build-all clean test fmt lint install run help

BINARY_NAME=inkwash
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
LDFLAGS=-ldflags "-X github.com/vexoa/inkwash/internal/config.Version=${VERSION} -X github.com/vexoa/inkwash/internal/config.Commit=${COMMIT} -X github.com/vexoa/inkwash/internal/config.Date=${DATE}"

help: ## Display this help message
	@echo "InkWash CLI - Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary for current platform
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME} ./cmd/inkwash

build-all: ## Build binaries for all platforms
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME}-linux-amd64 ./cmd/inkwash
	@GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME}-linux-arm64 ./cmd/inkwash
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME}-darwin-amd64 ./cmd/inkwash
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME}-darwin-arm64 ./cmd/inkwash
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME}-windows-amd64.exe ./cmd/inkwash

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf ${GOBIN}
	@go clean

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/"; \
		exit 1; \
	fi

install: build ## Install binary to system
	@echo "Installing ${BINARY_NAME}..."
	@sudo cp ${GOBIN}/${BINARY_NAME} /usr/local/bin/

run: build ## Build and run the CLI
	@${GOBIN}/${BINARY_NAME}

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy