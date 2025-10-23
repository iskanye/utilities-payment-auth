# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN apk add --no-cache build-base
RUN go build -o /go/bin/auth ./cmd/auth/main.go
RUN go build -o /go/bin/migrator ./cmd/migrator/main.go

# Запуск
FROM alpine
USER root
WORKDIR /home/app
COPY --from=builder /go/bin/auth ./
COPY --from=builder /go/bin/migrator ./
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
RUN mkdir storage
RUN ./migrator --storage-path=./storage/auth.db --migrations-path=./migrations
ENTRYPOINT ["./auth"]
CMD ["-config", "./config/dev.yaml"]