.DEFAULT_GOAL := help
.PHONY: help

build: ## Builds the application and all its dependencies as defined by the docker-compose.yml
	@docker-compose build

run: ## Runs the application and all its dependencies as defined by the docker-compose.yml
	@docker-compose up

test: ## Runs all tests, including integration ones
	@go test -race ./...

test-short: ## Runs all tests, excluding integration ones
	@go test -short -race ./...

vendor: ## Make vendored copy of dependencies
	@go mod vendor

.env:
	@cp .env.example .env

help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep ^help -v | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'
