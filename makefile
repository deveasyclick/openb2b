# Makefile for Tilvio project

# Variables
FRONTEND_DIR = frontend
BACKEND_DIR = api
MIGRATIONS_DIR=./api/db/migrations


# Default target
.PHONY: all
all: help

# Help message
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make help                 - Show this help message"
	@echo "  make install              - Install all dependencies"
	@echo "  make install-frontend     - Install frontend dependencies"
	@echo "  make install-backend      - Install backend dependencies"
	@echo "  make dev                  - Run both frontend and backend in development mode"
	@echo "  make frontend             - Run frontend in development mode"
	@echo "  make backend              - Run backend in development mode"
	@echo "  make build                - Build both frontend and backend"
	@echo "  make build-frontend       - Build frontend"
	@echo "  make build-backend        - Build backend"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make lint-frontend        - Lint frontend code"

.PHONY: install-backend
install-backend:
	@echo "Installing backend dependencies..."
	go mod tidy


.PHONY: backend
backend:
	@echo "Starting backend development server..."
	@cd $(BACKEND_DIR) && air

.PHONY: build-backend
build-backend:
	@echo "Building backend..."
	@cd ${BACKEND_DIR} && go build -o ./bin/main ./cmd/api/main.go

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(FRONTEND_DIR)/dist
	@# Add backend clean commands when implemented


.PHONY: lint-backend
lint-backend:
	@echo "Lint backend..."
	@cd backend && go fmt ./... && go vet ./...

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" status

migrate-new:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: swagger
swagger:
	cd api && swag init --output docs --generalInfo cmd/api/main.go
