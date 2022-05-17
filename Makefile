complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

up: ## Запуск проекта
	docker-compose up --build 

down: ## Остановка проета
	docker-compose down

build: ## Сборка проекта
	go build -o ./build/gophermart ./cmd/gophermart

exec: ## Зайти в контейнер
	docker-compose exec loyalty_backend sh

create_migrate: ## Создать миграцию
	migrate create -ext sql -dir db/migration -seq playground_schema 

migrate_up: ## Восстановить миграцию
	migrate -path db/migration -database "postgresql://loyalty:loyalty@localhost:5432/loyalty?sslmode=disable" -verbose up

migrate_down: ## Удалить миграцию
	migrate -path db/migration -database "postgresql://loyalty:loyalty@localhost:5432/loyalty?sslmode=disable" -verbose down