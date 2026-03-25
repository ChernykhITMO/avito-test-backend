APP_NAME ?= room-booking
ENV_FILE ?= .env.example

.PHONY: run build test fmt lint swagger-gen up down logs seed migrate-up migrate-down migrate-status smoke smoke-down

run:
	go run ./cmd/api

build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal ./tests

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run ./...

swagger-gen:
	go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/api/main.go -o ./docs/swagger --parseInternal

up:
	docker compose --env-file $(ENV_FILE) up --build

down:
	docker compose --env-file $(ENV_FILE) down --remove-orphans

logs:
	docker compose --env-file $(ENV_FILE) logs -f app

seed:
	@[ -n "$(DATABASE_URL)" ] || (echo "DATABASE_URL is required" && exit 1)
	go run ./cmd/seed

smoke:
	docker compose -f docker-compose.e2e.yaml --env-file $(ENV_FILE) up --build --abort-on-container-exit --exit-code-from smoke

smoke-down:
	docker compose -f docker-compose.e2e.yaml --env-file $(ENV_FILE) down -v --remove-orphans

migrate-up:
	@[ -n "$(DATABASE_URL)" ] || (echo "DATABASE_URL is required" && exit 1)
	goose -dir ./migrations postgres "$(DATABASE_URL)" up

migrate-down:
	@[ -n "$(DATABASE_URL)" ] || (echo "DATABASE_URL is required" && exit 1)
	goose -dir ./migrations postgres "$(DATABASE_URL)" down

migrate-status:
	@[ -n "$(DATABASE_URL)" ] || (echo "DATABASE_URL is required" && exit 1)
	goose -dir ./migrations postgres "$(DATABASE_URL)" status
