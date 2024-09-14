MAIN_PACKAGE_PATH := ./cmd
BINARY_NAME := main
MIGRATIONS_DIR := ./migrations
DB_DSN := "postgres://user:psswd@localhost/backend?sslmode=disable"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	export DB_DSN=$(DB_DSN)
	export JWT_SECRET="JWT_SECRET"
	set +a && source .env && set -a
	air \
		--build.cmd "go build -o /tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}" \
		--build.bin "/tmp/bin/main" \

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

## migrate/up: apply all up migrations
.PHONY: migrate/up
migrate/up:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) up

## migrate/down: apply all down migrations
.PHONY: migrate/down
migrate/down:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) down

## migrate/new: create a new migration
.PHONY: migrate/new
migrate/new:
	@[ "$(name)" ] || (echo "name is required" && exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
