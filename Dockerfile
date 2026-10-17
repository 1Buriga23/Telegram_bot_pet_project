
# Этап 1. Компилируем Go-приложение
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/gracia-bot ./cmd/bot

# Этап 2. Создаём лёгкий образ для запуска
FROM alpine:3.23

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/gracia-bot /app/gracia-bot

USER 10001:10001

ENTRYPOINT ["/app/gracia-bot"]
