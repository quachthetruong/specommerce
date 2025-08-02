# Makefile for Order Service

# Variables
BINARY_NAME=orderservice
MAIN_PATH=cmd/server/main.go
BUILD_DIR=build

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run
GOGENERATE=$(GOCMD) generate

# Swagger
SWAG_CMD=swag
SWAG_INIT=$(SWAG_CMD) init
SWAG_FMT=$(SWAG_CMD) fmt





.PHONY: all build clean test coverage deps generate-swagger generate-proto run dev help

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean completed!"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Tests completed!"

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed!"

# Generate Swagger documentation
generate-swagger:
	@echo "Generating Swagger documentation..."
	$(SWAG_INIT) -g $(MAIN_PATH) -o docs/openapi/api/orderservice
	@echo "Swagger documentation generated!"

# Generate protobuf code
generate-proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative model/*.proto
	@echo "Protobuf code generated!"

# Format Swagger documentation
format-swagger:
	@echo "Formatting Swagger documentation..."
	$(SWAG_FMT) -d docs/openapi/api/orderservice
	@echo "Swagger documentation formatted!"

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) $(MAIN_PATH)

# Run in development mode
dev:
	@echo "Running in development mode..."
	$(GORUN) $(MAIN_PATH)



# Install Swagger CLI
install-swagger:
	@echo "Installing Swagger CLI..."
	$(GOGET) -u github.com/swaggo/swag/cmd/swag
	@echo "Swagger CLI installed!"

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "Linting completed!"

# Install golangci-lint
install-lint:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
	@echo "golangci-lint installed!"

# Database migrations (configure in config.yml)
migrate-up:
	@echo "Running database migrations up..."
	@echo "Please configure database connection in config.yml"

migrate-down:
	@echo "Running database migrations down..."
	@echo "Please configure database connection in config.yml"

# Install migrate tool
install-migrate:
	@echo "Installing migrate tool..."
	$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate
	@echo "Migrate tool installed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run tests"
	@echo "  coverage           - Run tests with coverage"
	@echo "  deps               - Install dependencies"
	@echo "  generate-swagger   - Generate Swagger documentation"
	@echo "  generate-proto     - Generate protobuf code"
	@echo "  format-swagger     - Format Swagger documentation"
	@echo "  run                - Run the application"
	@echo "  dev                - Run in development mode"

	@echo "  install-swagger    - Install Swagger CLI"
	@echo "  lint               - Lint code"
	@echo "  install-lint       - Install golangci-lint"
	@echo "  migrate-up         - Run database migrations up"
	@echo "  migrate-down       - Run database migrations down"
	@echo "  install-migrate    - Install migrate tool"
	@echo "  help               - Show this help" 