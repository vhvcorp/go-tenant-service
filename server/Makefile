.PHONY: help build test lint clean run docker-build docker-push proto

# Variables
SERVICE_NAME := tenant-service
DOCKER_REGISTRY ?= ghcr.io/vhvplatform
VERSION ?= $(shell git describe --tags --always --dirty)
GO_VERSION := 1.25.5

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	@echo "Building services..."
	@go build -o bin/tenant-service ./cmd/tenant/main.go
	@go build -o bin/api-gateway ./cmd/gateway/main.go
	@echo "Build complete!"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ dist/ coverage.* *.out
	@go clean -testcache
	@echo "Clean complete!"

run-tenant: ## Run the tenant service locally
	@echo "Running tenant-service..."
	@go run ./cmd/tenant/main.go

run-gateway: ## Run the api-gateway service locally
	@echo "Running api-gateway..."
	@go run ./cmd/gateway/main.go

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

proto: ## Generate protobuf files
	@if [ -d "proto" ]; then \
		echo "Generating protobuf files..."; \
		protoc -I. -Iproto \
			--go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
			--openapiv2_out=docs --openapiv2_opt=logtostderr=true \
			proto/*.proto; \
	fi

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION) .
	@docker tag $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(SERVICE_NAME):latest
	@echo "Docker image built: $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION)"

docker-push: docker-build ## Push Docker image
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(SERVICE_NAME):latest
	@echo "Docker image pushed!"

docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 -p 50051:50051 \
		--name $(SERVICE_NAME) \
		$(DOCKER_REGISTRY)/$(SERVICE_NAME):latest

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@if [ -d "proto" ]; then \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest; \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest; \
	fi
	@echo "Tools installed!"

.DEFAULT_GOAL := help
