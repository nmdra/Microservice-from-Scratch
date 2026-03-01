FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-service main.go

FROM alpine:3.23

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/order-service .

EXPOSE 8080

CMD ["./order-service"]