# Load .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ---------- Variables ----------
GO             := go
MIGRATE_CMD    := migrate
MIGRATE_PATH   := db/migrations
SQLC_CMD       := sqlc
SQLC_CONFIG    := db/sqlc.yaml
GOLANGCI_LINT  := golangci-lint

# Database connection string
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# All targets are phony (not actual files)
.PHONY: help run tidy migrate-up migrate-down migrate-create migrate-version sqlc test test-race lint

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  run               Start the API server locally"
	@echo "  tidy              Clean and verify dependencies"
	@echo "  migrate-up        Apply all pending migrations"
	@echo "  migrate-down      Rollback the last migration"
	@echo "  migrate-create name=...   Create new migration"
	@echo "  migrate-version   Show current migration version"
	@echo "  sqlc              Generate Go code from sqlc queries"
	@echo "  test              Run all tests"
	@echo "  test-race         Run tests with race detector"
	@echo "  lint              Run golangci-lint"

# ---------- Development ----------
run:
	$(GO) run cmd/server/main.go

tidy:
	$(GO) mod tidy
	$(GO) mod verify

# ---------- Database ----------
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=<name>"; \
		exit 1; \
	fi
	$(MIGRATE_CMD) create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-up:
	$(MIGRATE_CMD) -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE_CMD) -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-version:
	$(MIGRATE_CMD) -path $(MIGRATE_PATH) -database "$(DB_URL)" version

# ---------- Tooling ----------
sqlc:
	$(SQLC_CMD) generate -f $(SQLC_CONFIG)

test:
	$(GO) test -v ./...

test-race:
	$(GO) test -race -v ./...

lint:
	$(GOLANGCI_LINT) run ./...

