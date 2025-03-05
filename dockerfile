# Build stage
FROM golang:alpine AS builder

# Установите рабочую директорию
WORKDIR /go/src/app

# Скопируйте зависимости
COPY go.mod go.sum ./
RUN go mod download

# Скопируйте код
COPY . .

# Устанавливаем переменные окружения для совместимости
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Собираем Go-приложения под Linux
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /go/bin/agent ./cmd/agent
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /go/bin/orchestrator ./cmd/orchestrator
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /go/bin/web ./web

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /go/bin/agent /app/agent
COPY --from=builder /go/bin/orchestrator /app/orchestrator
COPY --from=builder /go/bin/web /app/web

RUN chmod +x /app/agent /app/orchestrator /app/web

ENTRYPOINT ["/app/web"]

EXPOSE 8080
EXPOSE 8081
