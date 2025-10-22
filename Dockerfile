FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o auth .
FROM alpine
WORKDIR /build
COPY --from=builder /build/auth /build/auth
RUN chmod +x auth
CMD ["./auth", "--config=\"./config/prod.yaml\""]