# Server core

Это сервис, реализующий серверную логику для многопользовательской игры в морской бой. Отличительной особенностью является поддержка создания новых предметов на стороне сервиса инвентаря предоставленным метаязыком.

## Структура проекта

```text
.
├── cmd/                        # Точка входа (main.go)
├── docker-compose.yml          # Docker Compose для запуска приложения
├── Dockerfile                  # Сборка контейнера
├── grafana/                    # Дашборды и настройки мониторинга
├── internal/                   # Внутренняя логика приложения:
│   ├── api/                    #   - HTTP/Webhook роуты
│   ├── config/                 #   - инициализация и конфиг
│   ├── event/                  #   - события и диспетчеризация
│   ├── game/                   #   - игровая логика (ходы, стрельба, размещение кораблей)
│   ├── infra/                  #   - инфраструктурные сервисы (Kafka)
│   ├── items/                  #   - предметы и их скрипты
│   ├── match/                  #   - логика матчей и комнат
│   ├── transaction/            #   - транзакции и управление ими
│   ├── ws/                     #   - WebSocket обработчики
│   └── wsiface/                #   - интерфейсы для ws
├── packets/                    # Форматы пакетов и типов
├── prometheus.yml              # Конфиг для мониторинга
├── go.mod                      # Зависимости Go
├── go.sum                      # Контрольные суммы зависимостей
└── README.md                   # Документация
```

## Запуск

```bash
# Клонируйте репозиторий
git clone https://github.com/lesta-start-battleship/server-core.git

# Перейдите в папку с проектом
cd server-core

# генерация env файла
make env

# Запустите проект с помощью Docker
make build
```

После запуска приложение будет доступно по адресу: [http://localhost:8080](http://localhost8080)


## Зависимости

| Зависимость                  | Версия         |
|------------------------------|---------------|
| github.com/Shopify/sarama    | v1.38.2       |
| github.com/gorilla/websocket | v1.5.0        |
| github.com/prometheus/client_golang | v1.18.0 |
| github.com/joho/godotenv     | v1.5.1        |
| github.com/google/uuid       | v1.6.0        |
| github.com/stretchr/testify  | v1.8.4        |

(См. полный список в go.mod)

## Технологии

Go, Docker, Kafka, WebSocket, Prometheus, Grafana

## Мониторинг

Grafana:    http://localhost:3000/
Prometheus: http://localhost:9090/

## Тестирование
[Тут](docs/gameTestReport.md) описаны тест-кейсы для функций ядра игры.
