# Makefile for Inventory Service

# Variables
BINARY_NAME=inventory-service
PROTO_DIR=api/proto
GO_OUT_DIR=api/proto
MAIN_DIR=cmd

.PHONY: all proto build clean run test docker help

all: proto build

# Generate Go code from protobuf
proto:
	@echo "Generating Go files from protobuf..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "Proto files generated successfully"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_DIR)/main.go
	@echo "Binary built successfully"

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f bin/$(BINARY_NAME)
	@echo "Clean up completed"

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_DIR)/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests completed"

# Build docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .
	@echo "Docker image built successfully"

# Show help
help:
	@echo "Available commands:"
	@echo "  make proto    - Generate Go files from protobuf definitions"
	@echo "  make build    - Build the application"
	@echo "  make clean    - Clean up generated files"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make docker   - Build Docker image"
	@echo "  make help     - Show this help message"
