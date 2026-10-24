COMPOSE := docker compose
DB_DSN := postgres://postgres:test123@localhost:5433/habit-graph_app?sslmode=disable

.PHONY: help run-back run-front docker-up docker-down docker-logs migrate-up migrate-down test lint

help: ## Показать список всех команд
	@awk 'BEGIN {FS = ":.*##"; printf "\nCommands:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

run-back: ## Запустить Go-бэкенд локально
	go run ./cmd/main.go

run-front: ## Установить зависимости и запустить фронтенд локально (порт 3000)
	cd frontend && npm install && npm run dev

docker-up: ## Собрать и поднять всю инфраструктуру в Docker
	$(COMPOSE) up -d --build

docker-down: ## Остановить и удалить все контейнеры
	$(COMPOSE) down

docker-logs: ## Смотреть логи всех контейнеров в реальном времени
	$(COMPOSE) logs -f

migrate-create: ## Создать новую миграцию (использование: make migrate-create NAME=имя_миграции)
	@if [ -z "$(NAME)" ]; then echo "Ошибка: укажите имя (make migrate-create NAME=add_users)"; exit 1; fi
	docker run --rm -v $(CURDIR)/migrations:/migrations migrate/migrate:v4.18.1 \
		create -ext sql -dir /migrations -seq $(NAME)

migrate-up: ## Накатить миграции через контейнер migrate
	$(COMPOSE) run --rm migrate

migrate-down: ## Откатить последнюю миграцию
	docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate:v4.18.1 \
		-path=/migrations -database="$(DB_DSN)" down 1

db-reset: ## Сбросить базу данных 
	docker compose down -v
	docker compose up -d
	docker compose run --rm migrate


#soon
# test: ## Запустить unit-тесты Go с race-детектором
# 	go test -v -race -cover ./...

# lint: ## Запустить линтер golangci-lint
# 	golangci-lint run ./...