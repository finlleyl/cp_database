.PHONY: help build run test clean deps fmt lint docker-build docker-up docker-down docker-logs migrate-create migrate-up migrate-down

APP_NAME := cp_database
BINARY_NAME := app
CMD_PATH := ./cmd/app
DOCKER_COMPOSE := docker-compose --env-file .env -f infra/docker-compose.yml
DOCKERFILE := infra/Dockerfile

ifeq ($(OS),Windows_NT)
    DB_USER        ?= $(shell powershell -NoProfile -Command "(Get-Content .env | Where-Object {$$_ -match '^POSTGRES_USER='}) -replace '^POSTGRES_USER=',''")
    DB_PASSWORD    ?= $(shell powershell -NoProfile -Command "(Get-Content .env | Where-Object {$$_ -match '^POSTGRES_PASSWORD='}) -replace '^POSTGRES_PASSWORD=',''")
    DB_NAME        ?= $(shell powershell -NoProfile -Command "(Get-Content .env | Where-Object {$$_ -match '^POSTGRES_DB='}) -replace '^POSTGRES_DB=',''")
else
    DB_USER        ?= $(shell grep POSTGRES_USER .env | cut -d '=' -f2)
    DB_PASSWORD    ?= $(shell grep POSTGRES_PASSWORD .env | cut -d '=' -f2)
    DB_NAME        ?= $(shell grep POSTGRES_DB .env | cut -d '=' -f2)
endif
DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable

deps: ## Установить зависимости
	go mod download
	go mod tidy

build: ## Собрать приложение
	CGO_ENABLED=0 go build -o $(BINARY_NAME) $(CMD_PATH)

run: ## Запустить приложение локально
	go run $(CMD_PATH)

test: ## Запустить тесты
	go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt: ## Форматировать код
	go fmt ./...

lint:
	golangci-lint run ./...;

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -cache

docker-build: ## Собрать Docker образ
	docker build -f $(DOCKERFILE) -t $(APP_NAME):latest .

docker-up: ## Запустить сервисы через Docker Compose
	$(DOCKER_COMPOSE) up --build -d

docker-down: ## Остановить сервисы Docker Compose
	$(DOCKER_COMPOSE) down

docker-logs: ## Показать логи Docker Compose
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Перезапустить сервисы Docker Compose
	$(DOCKER_COMPOSE) restart

docker-ps: ## Показать статус контейнеров
	$(DOCKER_COMPOSE) ps

migrate-create: ## make migrate-create name=add_master_table
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path=migrations -database "$(DB_DSN)" up

migrate-down:
	migrate -path=migrations -database "$(DB_DSN)" down -all

swagger: ## Сгенерировать Swagger документацию
	swag init -g cmd/app/main.go -o docs

swagger-fmt: ## Отформатировать Swagger аннотации
	swag fmt