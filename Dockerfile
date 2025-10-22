# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go/bin/auth ./cmd/auth/main.go

# Запуск
FROM alpine
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --from=builder /go/bin/auth ./
COPY --from=builder /app/config ./config
ENTRYPOINT ["./auth"]
CMD ["-config", "./config/dev.yaml"]