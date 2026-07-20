# Microservice Wallet Core + Balances

Event-driven wallet ecosystem with:

- **Wallet Core** (Go) — producer: creates clients, accounts and transactions, publishes Kafka events
- **Balances** (Node.js / TypeScript / Fastify) — consumer: listens to `balance_updated`, persists balances and exposes a query API

## Architecture

```text
Client API
   |
   v
Wallet Core (Go :8080) ----publish----> Kafka topic: balance_updated
   |                                         |
   |                                         v
MySQL (wallet :3306)              Balances Service (Node :3003)
                                             |
                                             v
                                      MySQL (balances :3307)
```

When a transaction is created in Wallet Core, it publishes a `BalanceUpdated` event to Kafka. The Balances service consumes that event and upserts the current balance of both accounts.

## Requirements

- Docker
- Docker Compose

## Quick start

```bash
docker compose up -d --build
```

This single command starts:

| Service | Port | Description |
|---------|------|-------------|
| wallet-core | 8080 | Wallet Core API |
| balances | 3003 | Balances query API |
| mysql-wallet | 3306 | Wallet Core database |
| mysql-balances | 3307 | Balances database |
| kafka | 9092 | Message broker |
| zookeeper | 2181 | Kafka coordination |
| control-center | 9021 | Kafka UI |

Migrations and seed data run automatically on first database initialization. No manual scripts are required.

### Seeded accounts

| Account ID | Initial balance |
|------------|-----------------|
| `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa` | 1000 |
| `bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb` | 500 |

## Test the flow

1. Check balances (seeded values):

```bash
curl http://localhost:3003/balances/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
curl http://localhost:3003/balances/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb
```

2. Create a transaction in Wallet Core:

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_from_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "account_to_id": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    "amount": 50
  }'
```

3. Query balances again (should be 950 and 550):

```bash
curl http://localhost:3003/balances/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
curl http://localhost:3003/balances/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb
```

Ready-to-use HTTP requests are also available in [`api.http`](./api.http).

## Balances service (local development)

```bash
cd node
npm install
npm run dev
```

Environment defaults (override as needed):

| Variable | Default |
|----------|---------|
| `HTTP_PORT` | `3003` |
| `DB_HOST` | `localhost` |
| `DB_PORT` | `3307` |
| `DB_USER` | `root` |
| `DB_PASSWORD` | `root` |
| `DB_NAME` | `balances` |
| `KAFKA_BROKERS` | `localhost:9092` |
| `KAFKA_TOPIC_BALANCE_UPDATED` | `balance_updated` |

## Project structure

```text
.
├── cmd/wallet_core/          # Wallet Core entrypoint (Go)
├── internal/                 # Wallet Core domain, use cases, web, events
├── node/                     # Balances microservice (Node.js)
│   ├── database/             # MySQL schema + seeds
│   └── src/
│       ├── modules/balances/ # DDD module (domain, application, infra)
│       └── infra/            # HTTP, Kafka, MySQL wiring
├── scripts/mysql-wallet/     # Wallet Core schema + seeds
├── api.http                  # HTTP test collection
└── docker-compose.yaml
```

## Clean restart

If you need to recreate databases and re-run seeds:

```bash
docker compose down -v
docker compose up -d --build
```

## API

### `GET /balances/{account_id}`

Returns the current balance for an account.

```json
{
  "account_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
  "balance": 1000
}
```
