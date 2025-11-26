.DEFAULT_GOAL := help

.PHONY: help dev dev-down dev-logs build test migrate-up migrate-down sqlc web-dev web-build build-prod deploy-local

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Start all services with hot reload
	docker compose up --build

dev-down: ## Stop all services
	docker compose down

dev-logs: ## View logs
	docker compose logs -f

build: ## Build nexus binary
	go build -o bin/nexus .

test: ## Run tests
	go test ./...

migrate-up: ## Run migrations
	migrate -path internal/database/migrations -database "$$DATABASE_URL" up

migrate-down: ## Rollback migration
	migrate -path internal/database/migrations -database "$$DATABASE_URL" down 1

sqlc: ## Generate sqlc code
	sqlc generate

web-dev: ## Start frontend dev server
	cd web && bun run dev

web-build: ## Build frontend
	cd web && bun run build

# Production
build-prod: ## Build production Docker image
	docker build -f docker/Dockerfile.prod -t nexus:latest .

deploy-local: ## Test production deployment locally
	cd deploy && docker compose up -d --build
