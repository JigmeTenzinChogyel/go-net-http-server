# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Variables
MIGRATION_DIR=database/migrations
GOOSE_CMD=goose
DB_DSN=$(DB_CONN_STR)  # Replace with your actual database DSN

build:
	@go build -o bin/go-net-http-server cmd/main.go

test:
	@go test -v ./..

run: build
	@./bin/go-net-http-server

generate:
	@sqlc generate -f database/sqlc.yaml

# Create a new migration
create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create name=<migration_name>"; \
		exit 1; \
	fi; \
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) create $(name) sql

# Run migrations up
migrate-up:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) up

# Run migrations up by one
migrate-up-by-one:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) up-by-one

# Run migrations up to a specific version
migrate-up-to:
	@if [ -z "$(version)" ]; then \
		echo "Usage: make migrate-up-to version=<version>"; \
		exit 1; \
	fi; \
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) up-to $(version)

# Rollback migrations down
migrate-down:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) down

# Rollback migrations down to a specific version
migrate-down-to:
	@if [ -z "$(version)" ]; then \
		echo "Usage: make migrate-down-to version=<version>"; \
		exit 1; \
	fi; \
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) down-to $(version)

# Re-run the latest migration
migrate-redo:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) redo

# Rollback all migrations
migrate-reset:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) reset

# Status of migrations
migrate-status:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) status

# Print the current version of the database
migrate-version:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) version

# Fix sequential ordering of migrations
migrate-fix:
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) postgres $(DB_DSN) fix

full-reset: migrate-reset migrate-up generate