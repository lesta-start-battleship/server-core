# Server core

Это сервис, реализующий серверную логику для многопользовательской игры с поддержкой транзакций, событий, WebSocket и Kafka.

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

Для запуска необходим валидный `.env`

```bash
# Клонируйте репозиторий
git clone https://github.com/lesta-start-battleship/server-core.git

# Перейдите в папку с проектом
cd server-core

# Запустите проект с помощью Docker
docker-compose up --build -d
```

После запуска приложение будет доступно по адресу: [http://localhost:8080](http://localhost8080)

## Конфигурация

Указывается в `.env`

Пример `.env`:

```
# Port where game-core runs
GAME_CORE_PORT=8080

# Broker ips
KAFKA_BROKERS=37.9.53.228:9092

# Topics to send
USED_ITEMS=prod.inventory.fact.used-items.v1
MATCH_RESULTS=prod.game.fact.match-results.v1

# Base api url of invetary service
INVENTORY_SERVICE_GET_ALL_ITEMS=http://37.9.53.107/items/
INVENTORY_SERVICE_GET_USER_ITEMS=http://37.9.53.107/inventory/user_inventory
INVENTORY_SERVICE_USE_ITEM=http://37.9.53.107/inventory/use_item
```

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
