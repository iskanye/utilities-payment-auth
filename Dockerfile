# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/auth ./cmd/auth/main.go
RUN go build -o /bin/migrator ./cmd/migrator/main.go
RUN go build -o /bin/admins ./cmd/admins/main.go

# Запуск
FROM alpine
USER root
WORKDIR /home/app
COPY --from=builder /bin/auth ./
COPY --from=builder /bin/migrator ./
COPY --from=builder /bin/admins ./
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
RUN mkdir storage
RUN ./migrator --storage-path=./storage/auth.db --migrations-path=./migrations
RUN ./admins --storage=./storage/auth.db --admins=./config/admins.yaml
ENTRYPOINT ["./auth"]
CMD ["-config", "./config/dev.yaml"]