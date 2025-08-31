.PHONY: help build run dev clean docker-up docker-down docker-logs docker-clean deps lint format vet check install-tools

# Default target
help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Environment setup
.env: ## Create .env file from example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please edit .env with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Go commands
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

build: deps ## Build the application
	@echo "Building application..."
	@go build -o bin/api ./cmd/api

run: .env ## Run the application
	@echo "Starting API server..."
	@go run ./cmd/api

dev: .env install-tools ## Run in development mode with hot reload
	@echo "Starting development server with hot reload..."
	@air

# Code quality
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

check: format vet lint ## Run all checks (format, vet, lint)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -cache

# Docker commands
docker-up: ## Start MongoDB and services
	@echo "Starting Docker services..."
	@docker-compose up -d

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f

docker-restart: ## Restart Docker services
	@echo "Restarting Docker services..."
	@docker-compose restart

docker-clean: ## Stop services and remove volumes
	@echo "Cleaning Docker services and volumes..."
	@docker-compose down -v
	@docker-compose rm -f

# Database commands
db-connect: ## Connect to MongoDB shell
	@echo "Connecting to MongoDB..."
	@docker-compose exec mongodb mongosh -u admin -p password123 --authenticationDatabase admin user_management

db-backup: ## Backup database
	@echo "Creating database backup..."
	@mkdir -p backups
	@docker-compose exec mongodb mongodump --uri="mongodb://admin:password123@localhost:27017/user_management?authSource=admin" --out=/tmp/backup
	@docker-compose cp mongodb:/tmp/backup ./backups/backup-$(shell date +%Y%m%d-%H%M%S)
	@echo "Backup created in backups/ directory"

db-restore: ## Restore database (usage: make db-restore BACKUP=backup-20240101-120000)
	@if [ -z "$(BACKUP)" ]; then \
		echo "Usage: make db-restore BACKUP=backup-20240101-120000"; \
		exit 1; \
	fi
	@echo "Restoring database from $(BACKUP)..."
	@docker-compose cp ./backups/$(BACKUP) mongodb:/tmp/restore
	@docker-compose exec mongodb mongorestore --uri="mongodb://admin:password123@localhost:27017/user_management?authSource=admin" --drop /tmp/restore/user_management

# Development workflow
setup: install-tools .env docker-up ## Complete project setup
	@echo "Waiting for MongoDB to be ready..."
	@sleep 5
	@echo "Setup complete! You can now run 'make dev' to start development"

start: docker-up run ## Start database and API
	@echo "API is running at http://localhost:8080"

restart: docker-restart run ## Restart everything

# API testing
test-api: ## Test API endpoints (requires API to be running)
	@echo "Testing API endpoints..."
	@curl -s -o /dev/null -w "Health check: %{http_code}\n" http://localhost:8080/api/v1/health
	@echo "API endpoints tested"

# Production commands
build-prod: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t user-management-api .

# Show project status
status: ## Show project status
	@echo "=== Project Status ==="
	@echo "Go version: $(shell go version)"
	@echo "Docker status:"
	@docker-compose ps
	@echo "Git status:"
	@git status --porcelain || echo "Not a git repository"
