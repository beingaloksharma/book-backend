# Variables
APP_NAME := server
CMD_PATH := cmd/server/main.go
BUILD_DIR := bin

# Detect OS for binary extension
ifeq ($(OS),Windows_NT)
    BINARY_NAME := $(APP_NAME).exe
    RM_CMD := del /Q
    MKDIR_CMD := mkdir
else
    BINARY_NAME := $(APP_NAME)
    RM_CMD := rm -rf
    MKDIR_CMD := mkdir -p
endif

.PHONY: all build run test clean docs tidy deps help

all: build

## 🔨 Build: Compile the application
build:
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## 🚀 Run: Run the application
run:
	go run $(CMD_PATH)

## 🧪 Test: Run tests
test:
	go test -v ./...

## 🧪 Test: Run tests with coverage
test-cover:
	go test -cover ./...

## 📄 Docs: Generate Swagger documentation
docs:
	@echo "Generating Swagger documentation..."
	swag init -d cmd/server,internal/controller,internal/model -g main.go --parseDependency --parseInternal
	@echo "Done. Docs available at http://localhost:8080/swagger/index.html after running."

## 🧹 Clean: Cleanup build artifacts
clean:
	@echo "Cleaning..."
	-$(RM_CMD) $(BUILD_DIR)
	@echo "Clean complete."

## 📦 Deps: Download dependencies
deps:
	go mod download

## 🧹 Tidy: Format and tidy mod file
tidy:
	go fmt ./...
	go mod tidy

## ❓ Help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build   - Compile the application"
	@echo "  run     - Run the application"
	@echo "  test    - Run tests"
	@echo "  docs    - Generate Swagger documentation"
	@echo "  clean   - Remove build artifacts"
	@echo "  tidy    - Format code and tidy dependencies"
	@echo "  deps    - Download dependencies"
