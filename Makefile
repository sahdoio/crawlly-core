.PHONY: help go down logs rebuild clean status db-shell sh health run build dev dev-logs dev-down dev-rebuild

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

## Start with HOT RELOAD (recommended for development)
dev: down down-dev
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up
	@echo "🔥 Hot reload enabled! Edit files and see changes instantly."

dev-logs: ## View logs in dev mode
	docker compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

dev-down: ## Stop dev environment
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

dev-rebuild: ## Rebuild dev environment
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

## Start all services with Docker (production mode)
go: down down-dev
	docker compose up -d
	@echo "✅ Services started!"
	@echo "🌐 API: http://localhost:3000"
	@echo "📊 Health: http://localhost:3000/health"

down: ## Stop all services
	docker compose down
	@echo "✅ Services stopped!"

down-dev: ## Stop all services
	docker compose down
	@echo "✅ Services stopped!"

logs: ## View logs from all services
	docker compose logs -f

rebuild: ## Rebuild and restart all services
	docker compose up -d --build
	@echo "✅ Services rebuilt and restarted!"

clean: ## Stop services and remove all data
	docker compose down -v
	@echo "✅ All data cleaned!"

status: ## Show status of all containers
	docker compose ps

db-shell: ## Access PostgreSQL shell
	docker compose exec postgres psql -U postgres -d crawlly

sh: ## Access application container shell
	docker compose exec app sh

health: ## Test health endpoint
	@curl http://localhost:3000/health
	@echo ""

run: ## Run application locally (without Docker)
	go run cmd/api/main.go

build: ## Build application binary
	go build -o bin/crawlly cmd/api/main.go
	@echo "✅ Binary created at bin/crawlly"
