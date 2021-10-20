# Include variables from the .envrc file
include .envrc

# ======================================================================== #
# HELPERS
# ======================================================================== #

## help: print this help message.
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [Y/N] ' && read ans && [ $${ans:-N} = y ]

# ======================================================================== #
# CODE QUALITY CONTROL
# ======================================================================== #

## audit: perform all code quality checks and module dependency resolution & verification.
.PHONY: audit
audit: format
	@echo '.....Tidying, verifying and resolving module dependencies.....'
	go mod tidy
	go mod verify

.PHONY: fmt
fmt:
	@echo '.....Formatting go code.....'
	go fmt ./...

.PHONY:lint
lint: fmt
	@echo '.....Linting go code.....'
	golint ./...

.PHONY: vet
vet: lint
	@echo '.....Vetting go code.....'
	go vet ./...

## format: format & lint all go files.
.PHONY: format
format: vet

.PHONY: vendor
vendor: audit
	@echo '.....Vendoring app dependencies.....'
	go mod vendor

# ======================================================================== #
# DEVELOPMENT
# ======================================================================== #

## start: run the ./cmd/api/ application.
.PHONY: start
start: check-cdb-env
	@go run ./cmd/api/ -db-dsn=${MONGODB_DSN}

## dsn_missing_error: Error string returned in an event where the CDB_DSN env var is missing / undefined
define dsn_missing_error

MONGODB_DSN envvar is undefined. To start the server, this envvar must be provided.
Create the envvar by running: 

export MONGODB_DSN=mongodb://localhost:27017/shopit

You can replace the above example URI with either your own local mongoDB setup URI or the mongo atlas remote URI.

endef

export dsn_missing_error

## check-MONGODB_DSN-env: Checks for the availability of the [MONGODB_DSN] env var and returns the [dsn_missing_error] if undefined
.PHONY: check-cdb-env
check-cdb-env:
ifndef MONGODB_DSN
		$(error ${dsn_missing_error})
endif