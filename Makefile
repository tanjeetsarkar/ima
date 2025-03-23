# Variables
APP_NAME := image-archive
SRC_DIR := .
BUILD_DIR := build
GO_FILES := $(shell find $(SRC_DIR) -name '*.go')
OS := $(shell uname -s)

# Default target
.PHONY: all
all: build

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(SRC_DIR)/main.go

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows.exe $(SRC_DIR)/main.go

# Build for the current OS
.PHONY: build
build:
	@echo "Building for current OS ($(OS))..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

# Run the application
.PHONY: run
run:
	@echo "Running the application..."
	go run $(SRC_DIR)/main.go

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all            - Build the project for the current OS"
	@echo "  build          - Build the project for the current OS"
	@echo "  build-linux    - Build the project for Linux"
	@echo "  build-windows  - Build the project for Windows"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  help           - Show this help message"