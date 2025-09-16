.PHONY: build test migrate-up migrate-down docker-up docker-down clean clean-all

# Build binaries
build:
	go build -o bin/consumer ./cmd/consumer
	go build -o bin/server ./cmd/server

# Run tests
test:
	go test ./internal/cache/... -v
	go test ./internal/repository/... -v
	go test ./internal/service/... -v

# Database migrations
migrate-up:
	goose -dir migrations postgres "postgres://app_user:secure_password_123@localhost:5432/orders?sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "postgres://app_user:secure_password_123@localhost:5432/orders?sslmode=disable" down

# Docker commands
docker-up:
	docker-compose -f deployments/docker-compose.yml up -d --build

docker-down:
	docker-compose -f deployments/docker-compose.yml down

docker-build:
	docker build -t l0-consumer -f consumer.Dockerfile .
	docker build -t l0-server -f server.Dockerfile .

# Clean up (Windows compatible)
clean:
	@echo "Cleaning up..."
	@if exist bin rmdir /s /q bin
	@docker-compose -f deployments/docker-compose.yml down -v
	@echo "Clean complete"

# Full clean including Docker system
clean-all: clean
	@echo "Performing full system cleanup..."
	@docker system prune -a -f --volumes
	@echo "Full cleanup complete"

# Development (run without docker)
run-consumer:
	go run cmd/consumer/main.go

run-server:
	go run cmd/server/main.go

# Setup (migrate and run)
setup: migrate-up docker-up

# Help
help:
	@echo "Available commands:"
	@echo "  make build        - Build binaries"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make clean        - Clean project (Windows compatible)"
	@echo "  make clean-all    - Full system cleanup"
	@echo "  make run-consumer - Run consumer locally"
	@echo "  make run-server   - Run server locally"
	@echo "  make setup        - Setup project (migrate + start containers)"