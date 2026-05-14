FROM golang:1.26 AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /time4book ./cmd/api

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /time4book /time4book

EXPOSE 50052

ENTRYPOINT ["/time4book"]
