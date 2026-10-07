FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./
COPY --from=builder /app/internal/app/adapters/postgres/sql ./internal/app/adapters/postgres/sql

RUN touch .env

ENV ENV=prod \
    GOOSE_MIGRATION_DIR=internal/app/adapters/postgres/sql \
    GOOSE_DRIVER=postgres

EXPOSE 50052

CMD ["./main"]
