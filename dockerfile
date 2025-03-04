# Используем официальный образ Go для сборки
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код в контейнер
COPY . .

# Сборка приложения
RUN go build -o agent ./cmd/agent/main.go
RUN go build -o orchestrator ./cmd/orchestrator/main.go
RUN go build -o web ./web/main.go

# Используем более легкий образ для запуска
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates

# Копируем собранные бинарники из этапа сборки
COPY --from=builder /app/agent /usr/local/bin/agent
COPY --from=builder /app/orchestrator /usr/local/bin/orchestrator
COPY --from=builder /app/web /usr/local/bin/web

# Копируем статические файлы
COPY --from=builder /app/web/static /usr/local/bin/static
COPY --from=builder /app/web/templates /usr/local/bin/templates

# Указываем команду по умолчанию для запуска веб-приложения
CMD ["web"]
