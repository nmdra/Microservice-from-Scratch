# Microservice-from-Scratch

A simple microservices demo using three services written in different languages, backed by PostgreSQL and fronted by an Nginx API gateway — all orchestrated with Docker Compose.

## Services

| Service | Language / Framework | Route prefix |
|---|---|---|
| item-service | Node.js / Express | `/items` |
| order-service | Go / Gin | `/orders` |
| payment-service | Python / Flask | `/payments` |

## Architecture

```mermaid
architecture-beta
    group api(cloud)[API Gateway]
    group services(cloud)[Microservices]
    group db(database)[Database]

    service client(internet)[Client]
    service nginx(server)[NGINX] in api
    service item(server)[Item Service<br/>Node.js] in services
    service order(server)[Order Service<br/>Go] in services
    service payment(server)[Payment Service<br/>Python] in services
    service postgres(database)[PostgreSQL] in db

    client:R --> L:nginx
    nginx:R --> L:item
    nginx:R --> L:order
    nginx:R --> L:payment
    item:B --> T:postgres
    order:B --> T:postgres
    payment:B --> T:postgres
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