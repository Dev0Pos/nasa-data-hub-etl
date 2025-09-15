# NASA Data Hub ETL Makefile

# Variables
APP_NAME := nasa-data-hub-etl
VERSION := 1.0.0
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_IMAGE_LATEST := $(APP_NAME):latest

# Go variables
GO_VERSION := 1.24
GOOS := linux
GOARCH := amd64

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help build test clean docker-build docker-run run lint fmt vet

# Default target
help: ## Show this help message
	@echo "$(BLUE)NASA Data Hub ETL - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development targets
build: ## Build the application
	@echo "$(BLUE)Building $(APP_NAME)...$(NC)"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo -o bin/$(APP_NAME) ./cmd/etl
	@echo "$(GREEN)Build completed successfully!$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean completed!$(NC)"

# Docker targets
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE_LATEST)
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(NC)"

docker-run: ## Run Docker container locally
	@echo "$(BLUE)Running Docker container...$(NC)"
	@docker run --rm -p 8080:8080 -v $(PWD)/config.yaml:/app/config.yaml:ro $(DOCKER_IMAGE_LATEST)

# Development targets
run: ## Run the application locally
	@echo "$(BLUE)Running application locally...$(NC)"
	@go run cmd/etl/main.go


# Utility targets
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy

deps-update: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy

version: ## Show version information
	@echo "$(BLUE)Version Information:$(NC)"
	@echo "Application: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "Docker Image: $(DOCKER_IMAGE)"

# All-in-one targets
dev-setup: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@make deps
	@make fmt
	@make vet
	@echo "$(GREEN)Development environment ready!$(NC)"
