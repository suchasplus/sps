# Go project variables
BINARY_NAME=sps
CMD_PATH=./cmd/sps

# Go commands
GO_BUILD=go build
GO_TEST=go test
GO_CLEAN=go clean
GO_LINT=golangci-lint

.PHONY: all build run test clean lint deps help

all: build ## (Default) Builds the application binary.

build: ## Builds the application binary.
	@echo "Building $(BINARY_NAME)..."
	@$(GO_BUILD) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "$(BINARY_NAME) built successfully."

run: build ## Runs the application. Pass arguments using "make run ARGS='your-args-here'"
	@./$(BINARY_NAME) $(ARGS)

test: ## Runs all unit tests with verbose output.
	@echo "Running tests..."
	@$(GO_TEST) -v ./...

clean: ## Removes the compiled binary.
	@echo "Cleaning up..."
	@$(GO_CLEAN)
	@rm -f $(BINARY_NAME)
	@echo "Cleanup complete."

lint: ## Lints the Go source code using golangci-lint.
	@echo "Linting code..."
	@if ! command -v $(GO_LINT) &> /dev/null; then \
		echo "golangci-lint not found. Please run 'make deps'."; \
		exit 1; \
	fi
	@$(GO_LINT) run

deps: ## Installs development dependencies (linter).
	@echo "Installing dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

help: ## Displays this help message.
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

