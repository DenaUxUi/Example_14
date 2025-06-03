# Этап 1: Сборка Go-приложения
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем Go-файлы
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем исполняемый файл
RUN go build -o server main.go

# Этап 2: Финальный образ
FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/server .

# Копируем статические файлы (если они нужны)
COPY static/ ./static

# Указываем порт, который слушает Go-сервер
EXPOSE 8080

# Запуск Go-сервера
CMD ["./server"]

