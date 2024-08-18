# Используем официальный образ Go для сборки
FROM golang:1.20-alpine AS builder

# Создаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файлы go.mod и go.sum, если он уже есть
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN go build -o bot

# Минимальный образ для запуска
FROM alpine:3.18

# Создаем рабочую директорию для приложения
WORKDIR /app

# Копируем скомпилированное приложение из стадии сборки
COPY --from=builder /app/bot .

# Запуск приложения
CMD ["./bot"]
