.PHONY: all build test clean run coverage fmt lint install help

# Binary name
BINARY_NAME=multi-agent-graph
# Main package path
MAIN_PATH=./examples/cmd/main.go
# Build output directory
BUILD_DIR=./build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

all: build

## build: Build the CLI binary
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

## install: Build and install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Install complete: $(GOPATH)/bin/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	@$(GOTEST) ./... -v

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) ./... -coverprofile=coverage.out
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-coverage-short: Run tests with coverage (no HTML)
test-coverage-short:
	@echo "Running tests with coverage..."
	@$(GOTEST) ./... -cover

## test-unit: Run unit tests only (skip integration tests)
test-unit:
	@echo "Running unit tests..."
	@$(GOTEST) -short ./... -v

## clean: Remove build artifacts and coverage files
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out
	@rm -f coverage.html
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -s -w .
	@echo "Format complete"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@$(GOLINT) run ./...
	@echo "Lint complete"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@echo "Dependencies downloaded"

## deps-tidy: Tidy and verify dependencies
deps-tidy:
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy
	@$(GOMOD) verify
	@echo "Dependencies tidied"

## run-linear: Run linear workflow
run-linear: build
	@echo "Running linear workflow..."
	@./$(BINARY_NAME) linear

## run-linear-parallel: Run linear workflow in parallel mode
run-linear-parallel: build
	@echo "Running linear workflow (parallel)..."
	@./$(BINARY_NAME) linear --parallel

## run-parallel: Run parallel workflow
run-parallel: build
	@echo "Running parallel workflow..."
	@./$(BINARY_NAME) parallel --parallel

## run-complex: Run complex workflow
run-complex: build
	@echo "Running complex workflow..."
	@./$(BINARY_NAME) complex --parallel

## run-etl: Run ETL pipeline
run-etl: build
	@echo "Running ETL pipeline..."
	@./$(BINARY_NAME) etl --parallel

## run-cicd: Run CI/CD pipeline
run-cicd: build
	@echo "Running CI/CD pipeline..."
	@./$(BINARY_NAME) cicd --parallel

## run-ml: Run ML pipeline
run-ml: build
	@echo "Running ML pipeline..."
	@./$(BINARY_NAME) ml --parallel

## run-book-flight: Run flight booking workflow
run-book-flight: build
	@echo "Running flight booking workflow..."
	@./$(BINARY_NAME) book-flight --parallel

## run-all: Run all example workflows
run-all: build
	@echo "Running all workflows..."
	@./$(BINARY_NAME) linear --parallel
	@./$(BINARY_NAME) parallel --parallel
	@./$(BINARY_NAME) complex --parallel
	@./$(BINARY_NAME) etl --parallel
	@./$(BINARY_NAME) cicd --parallel
	@./$(BINARY_NAME) ml --parallel
	@./$(BINARY_NAME) book-flight --parallel

test-nexus:
	cd nexus && go test -v ./...

build-nexus:
	cd nexus && go build ./...

## ci: Run CI checks (fmt, test, lint)
ci: fmt test
	@echo "CI checks passed"

## help: Display this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
