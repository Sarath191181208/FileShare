MAIN_PACKAGE_PATH := ./cmd
BINARY_NAME := main

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
	export DB_DSN="postgres://user:psswd@localhost/backend?sslmode=disable" 
	export JWT_SECRET="JWT_SECRET"
	set +a && source .env && set -a
	air \
		--build.cmd "go build -o /tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}" \
		--build.bin "/tmp/bin/main" \
