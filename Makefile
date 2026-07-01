# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

BINARY_NAME=api
BUILD_DIR=bin
MAIN_PACKAGE=./cmd/api

# -----------------------------------------------------------------------------
# Help
# -----------------------------------------------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "Build:"
	@echo "  make build          Build the application"
	@echo "  make run            Run the application locally"
	@echo "  make clean          Remove compiled binaries"
	@echo ""
	@echo "Development:"
	@echo "  make fmt            Format Go source files"
	@echo "  make vet            Run go vet"
	@echo "  make test           Run all tests"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   Build Docker images"
	@echo "  make up             Start database, migrate, seed and API"
	@echo "  make down           Stop containers"
	@echo "  make destroy        Stop containers and remove volumes"
	@echo "  make restart        Restart the application"
	@echo "  make logs           Show container logs"
	@echo "  make ps             Show running containers"

# -----------------------------------------------------------------------------
# Build
# -----------------------------------------------------------------------------

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# -----------------------------------------------------------------------------
# Run
# -----------------------------------------------------------------------------

.PHONY: run
run:
	go run $(MAIN_PACKAGE)

# -----------------------------------------------------------------------------
# Clean
# -----------------------------------------------------------------------------

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# -----------------------------------------------------------------------------
# Format
# -----------------------------------------------------------------------------

.PHONY: fmt
fmt:
	go fmt ./...

# -----------------------------------------------------------------------------
# Vet
# -----------------------------------------------------------------------------

.PHONY: vet
vet:
	go vet ./...

# -----------------------------------------------------------------------------
# Tests
# -----------------------------------------------------------------------------

.PHONY: test
test:
	go test ./...

# -----------------------------------------------------------------------------
# Docker
# -----------------------------------------------------------------------------

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: up-db
up-db:
	docker compose up -d db
	@echo "Waiting for PostgreSQL..."
	@until docker compose exec -T db pg_isready -U restaurant -d restaurant; do \
		sleep 1; \
	done

.PHONY: migrate
migrate:
	docker compose --profile migrate run --rm migrate

.PHONY: seed
seed:
	docker compose exec -T db psql -U restaurant -d restaurant < scripts/seed.sql

.PHONY: up-api
up-api:
	docker compose up -d api

.PHONY: up
up: up-db migrate seed up-api

.PHONY: down
down:
	docker compose down

.PHONY: destroy
destroy:
	docker compose down -v

.PHONY: restart
restart:
	$(MAKE) down
	$(MAKE) up

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: ps
ps:
	docker compose ps

.PHONY: shell-db
shell-db:
	docker compose exec db psql -U restaurant -d restaurant
