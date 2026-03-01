# Microservice-from-Scratch

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