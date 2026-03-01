# Microservice-from-Scratch

A simple microservices demo using three services written in different languages, backed by PostgreSQL and fronted by an Nginx API gateway — all orchestrated with Docker Compose.

## Services

| Service | Language / Framework | Route prefix |
|---|---|---|
| item-service | Node.js / Express | `/items` |
| order-service | Go / Gin | `/orders` |
| payment-service | Python / Flask | `/payments` |

## Architecture

```
Client → Nginx (port 8080) → item-service
                           → order-service
                           → payment-service
                           ↕
                        PostgreSQL
```

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose

## Quick Start

```bash
# 1. Copy and fill in environment variables
cp .env.example .env

# 2. Start all services
docker compose up --build

# 3. Check gateway health
curl http://localhost:8080/health
```

## License

[MIT](LICENSE)
