.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest

build: ## Build docker image to deploy
	docker build -t radish-miyazaki/go-web-app:${DOCKER_TAG} --target deploy ./

build-local: ## Build docker image to local development
	docker compose build --no-cache

up: ## Do docker compose up with hot reload
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

test: ## Execute tests
	go test -race -shuffle=on ./...

migrate: ## Read schema SQL from the file and execute migrate
	mysqldef -u todo -p password -h 127.0.0.1 -P 33060 todo < ./_tools/mysql/schema.sql

dry-migrate: ## Execute migrate with dry-run option
	mysqldef -u todo -p password -h 127.0.0.1 -P 33060 todo --dry-run < ./_tools/mysql/schema.sql

generate: ## Generate codes
	go generate ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
