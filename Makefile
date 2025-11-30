.PHONY: help build run test clean deps fmt lint docker-build docker-up docker-down docker-logs

APP_NAME := cp_database
BINARY_NAME := app
CMD_PATH := ./cmd/app
DOCKER_COMPOSE := docker-compose --env-file .env -f infra/docker-compose.yml
DOCKERFILE := infra/Dockerfile

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
