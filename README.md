# Universe Group — Golang Developer Test Task

Two microservices: **Products** (CRUD + event publishing) and **Notifications** (event consumer + logging).

## Tech Stack

- **Go 1.25** — Chi router, pgx (raw SQL), segmentio/kafka-go
- **PostgreSQL 16** — database
- **Apache Kafka** — message broker (KRaft mode)
- **Prometheus** — metrics
- **golang-migrate** — database migrations

## Project Structure

```
├── cmd/
│   ├── products/main.go          # Products service entry point
│   └── notifications/main.go     # Notifications service entry point
├── internal/products/             # Business logic (handler, service, repository, model)
├── pkg/kafka/                     # Shared Kafka producer & consumer
├── migrations/                    # SQL migrations
├── docker-compose.yml             # PostgreSQL, Kafka, Prometheus
└── prometheus.yml                 # Prometheus config
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### 1. Start infrastructure

```bash
docker-compose up -d
```

This starts PostgreSQL (port 5432), Kafka (port 9092), and Prometheus (port 9090).

### 2. Run Products service

```bash
go run ./cmd/products/
```

The service starts on `:8080`. Migrations are applied automatically.

### 3. Run Notifications service

In a separate terminal:

```bash
go run ./cmd/notifications/
```

Listens for product events from Kafka and logs them.

## API Endpoints

### Create product

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "iPhone 16", "description": "Smartphone", "price": 999.99}'
```

**Response:** `201 Created`

### List products (with pagination)

```bash
curl "http://localhost:8080/products?page=1&page_size=10"
```

**Response:** `200 OK`

### Delete product

```bash
curl -X DELETE http://localhost:8080/products/1
```

**Response:** `204 No Content`

## Running Tests

```bash
go test ./... -v
```

## Prometheus Metrics

Available at `http://localhost:8080/metrics`.

- `products_created_total` — total products created
- `products_deleted_total` — total products deleted

Prometheus UI: `http://localhost:9090`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://postgres:root@localhost:5432/products_db?sslmode=disable` | PostgreSQL connection string |
| `HTTP_ADDR` | `:8080` | Products HTTP server address |
| `MIGRATIONS_PATH` | `file://migrations` | Path to migration files |
| `KAFKA_BROKER` | `localhost:9092` | Kafka broker address |
| `KAFKA_TOPIC` | `product-events` | Kafka topic name |
| `KAFKA_GROUP` | `notifications-group` | Kafka consumer group (Notifications) |
