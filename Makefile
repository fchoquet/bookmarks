# Env Variables
# =============================================================================================
ROOT_DIR=$(PWD)
BINARY_PATH=$(ROOT_DIR)/bin/bookmarks

# Rules
# =============================================================================================
.PHONY: usage init build up test
.DEFAULT: usage

usage:
	@echo '+---------------------------------------------------------------------------------------+'
	@echo '| Bookmarks "make" Usage                                                                  |'
	@echo '+---------------------------------------------------------------------------------------+'
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "|- \033[33m%-15s\033[0m -> %s\n", $$1, $$2}'

init: ## Initializes the dev environment (WARNING: will overwrite DB)
	@echo 'initializing DB' && ./database/provision.sh

build: ## Builds the go binary
	@echo "Compiling..." && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_PATH) $(ROOT_DIR)/main.go

up: build ## All-in-one command that get up to date containers up and running
	@docker-compose build && docker-compose up -d mysql && docker-compose up -d api

test: ## Launches the unit test suite
	@go test -v $$(go list ./... | grep -v /vendor/)
