# Examples:
# make build-image tag=my-service:1.0.0
.PHONY: build-image
build-image:
	@./scripts/build.sh $(tag)

# Examples:
# make generate-wire
.PHONY: generate-wire
generate-wire:
	@./scripts/wire.sh

# Examples:
# make run service=auth
.PHONY: run
run:
	@echo "Run app..."
	@go run ./cmd/app

# Examples:
# make create-migration name=create_products_table
.PHONY: create-migration
create-migration:
	@echo "Create migration..."
	@migrate create -ext sql -dir pkg/db/migrations -seq $(name)

# Examples:
# make migrate-up dsn='postgres://user:password@localhost:5432/db?sslmode=disable' step=1
.PHONY: migrate-up
migrate-up:
	@echo "Migrate up..."
	@migrate -path pkg/db/migrations -database $(dsn) -verbose up $(step)

# Examples:
# make migrate-down dsn='postgres://user:password@localhost:5432/db?sslmode=disable' step=1
.PHONY: migrate-down
migrate-down:
	@echo "Migrate down..."
	@migrate -path pkg/db/migrations -database $(dsn) -verbose down $(step)

# Examples:
# make generate-swagger
.PHONY: generate-swagger
generate-swagger:
	@./scripts/doc.sh

# Example:
# make lint
.PHONY: lint
lint:
	@./scripts/lint.sh
