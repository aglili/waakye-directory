# Load environment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

# Variables
PROJECT_NAME := waakye-directory
MIGRATION_DIR := migrations
DB_URL := "postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5433/$(DB_NAME)?sslmode=disable"
DOCKER_COMPOSE := docker-compose
GO := go

# Go Tooling
.PHONY: all fmt lint test build run clean

all: fmt lint test build

fmt:
	@echo "Running go fmt..."
	@$(GO) fmt ./...

lint:
	@echo "Running go vet..."
	@$(GO) vet ./...

test:
	@echo "Running tests..."
	@$(GO) test ./... -cover

build:
	@echo "Building the project..."
	@$(GO) build -o bin/$(PROJECT_NAME) ./cmd/main.go

run:
	@echo "Running the application..."
	@./bin/$(PROJECT_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf bin/*

# Docker Commands
.PHONY: docker-up docker-down docker-restart docker-logs

docker-up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d --build

docker-down:
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

docker-restart:
	@echo "Restarting Docker containers..."
	@$(DOCKER_COMPOSE) down && $(DOCKER_COMPOSE) up -d --build

docker-logs:
	@echo "Showing Docker logs..."
	@$(DOCKER_COMPOSE) logs -f

# Database Migrations
.PHONY: migrate-create migrate-up migrate-down migrate-force

migrate-create:
	@echo "Creating new migration..."
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-up:
	@echo "Applying migrations with DB_URL=$(DB_URL)"
	@migrate -path $(MIGRATION_DIR) -database $(DB_URL) up

migrate-down:
	@echo "Reverting migrations..."
	@migrate -path $(MIGRATION_DIR) -database $(DB_URL) down

migrate-force:
	@echo "Forcing migration version..."
	@migrate -path $(MIGRATION_DIR) -database $(DB_URL) force $(version)

# Utility
.PHONY: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  fmt                Run go fmt"
	@echo "  lint               Run go vet"
	@echo "  test               Run tests with coverage"
	@echo "  build              Build the Go project"
	@echo "  run                Run the built executable"
	@echo "  clean              Remove built artifacts"
	@echo ""
	@echo "  docker-up          Start Docker containers"
	@echo "  docker-down        Stop Docker containers"
	@echo "  docker-restart     Restart Docker containers"
	@echo "  docker-logs        View Docker logs"
	@echo ""
	@echo "  migrate-create     Create a new database migration (use name=your_migration)"
	@echo "  migrate-up         Run all up migrations"
	@echo "  migrate-down       Run all down migrations"
	@echo "  migrate-force      Force a specific migration version (use version=version_number)"
