# Используем официальный образ Go для сборки
FROM golang:1.22 AS builder

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

CMD ["./bot"]