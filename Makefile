include .env
export

.PHONY: help build run dev migrate-up migrate-down migrate-create clean docker-up docker-down test

help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application"
	@echo "  make dev                - Run with air (live reload)"
	@echo "  make migrate-up         - Apply all pending migrations"
	@echo "  make migrate-down       - Revert last migration"
	@echo "  make migrate-create     - Create a new migration file (use NAME=migration_name)"
	@echo "  make sqlc-gen           - Generate type-safe database code from SQL"
	@echo "  make docker-up          - Start Docker containers (PostgreSQL + API)"
	@echo "  make docker-down        - Stop Docker containers"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make test               - Run tests"

# Build
build:
	go build -o .build/rest-api ./cmd/rest

# Run
run: build
	./.build/rest-api

# Development with live reload
dev:
	air

# Migrations
migrate-up:
	migrate -path ./internal/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	migrate -path ./internal/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./internal/database/migrations -seq $(NAME)

# SQLC
sqlc-gen:
	@echo "Generating type-safe database code..."
	sqlc generate
	@echo "✓ Code generated!"

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Testing
test:
	go test -v -cover ./...

# Clean
clean:
	rm -f rest-api
	go clean
	rm -rf tmp/

# Format and lint
fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Setup (for first time)
setup: deps
	@echo "Installing migrate CLI..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing air..."
	go install github.com/cosmtrek/air@latest
	@echo "Installing sqlc..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "✓ Setup complete!"

# Generate all migrations and run them
setup-db: docker-up
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 5
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up
