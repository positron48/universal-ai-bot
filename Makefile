APP_NAME := tgbot-skeleton
GO := go
DOCKER_IMAGE := tgbot-skeleton
DOCKER_TAG := latest

-include .env
.EXPORT_ALL_VARIABLES:

.PHONY: all tidy build run test lint setup up clean

all: build

# Go commands
tidy:
	$(GO) mod tidy

build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/bot

run: tidy build
	./bin/$(APP_NAME)

test:
	$(GO) test ./...

test-verbose:
	$(GO) test -v ./...

# Linting
GOLANGCI := $(shell if [ -x ./bin/golangci-lint ]; then echo ./bin/golangci-lint; else echo golangci-lint; fi)

lint:
	$(GOLANGCI) run --timeout=3m

lint-install:
	@echo "Installing golangci-lint v1.61.0 into ./bin..."
	@mkdir -p bin
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.61.0

# Setup
setup:
	@echo "Setting up project..."
	@mkdir -p bin
	@cp env.example .env
	@echo "✅ Project setup complete!"
	@echo "📝 Please edit .env file with your bot token"

# Development
dev: tidy
	$(GO) run ./cmd/bot

up: run

# Cleanup
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker compose up -d

docker-stop:
	docker compose down

docker-logs:
	docker compose logs -f

docker-clean:
	docker compose down -v
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true

docker-rebuild: docker-clean docker-build docker-run

# Development with Docker
docker-dev:
	docker compose up -d

docker-dev-logs:
	docker compose logs -f

docker-dev-restart:
	docker compose restart tgbot-skeleton

# Deployment commands
deploy:
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh deploy

docker-deploy: docker-build docker-run

update:
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh update

status:
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh status

logs:
	@chmod +x scripts/deploy.sh && ./scripts/deploy.sh logs

# Help
help:
	@echo "Available commands:"
	@echo "  make setup          - Initial project setup"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run in development mode"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linter"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run with docker compose"
	@echo "  make docker-stop    - Stop docker compose"
	@echo "  make docker-logs    - Show docker logs"
	@echo "  make docker-clean   - Clean Docker resources"
	@echo "  make docker-deploy  - Deploy with Docker"
	@echo ""
	@echo "Deployment commands:"
	@echo "  make deploy         - Deploy binary from GitHub releases"
	@echo "  make update         - Update deployed binary"
	@echo "  make status         - Check service status"
	@echo "  make logs           - Show service logs"

