.PHONY: env up

# Создать .env файл из примера, если он не существует
env:
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			cp .env.example .env; \
			echo "Файл .env создан из .env.example"; \
		else \
			echo "Ошибка: файл .env.example не найден"; \
			exit 1; \
		fi \
	else \
		echo "Файл .env уже существует"; \
	fi

# Запустить docker-compose
up:
	docker-compose up -d

# Остановить docker-compose
down:
	docker-compose down

# Перезапустить docker-compose
restart: down up

# Показать логи docker-compose
logs:
	docker-compose logs -f

# Собрать и запустить docker-compose
build:
	docker-compose up -d --build

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./... && \
	go tool cover -html=coverage.out
