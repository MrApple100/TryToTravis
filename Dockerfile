# Этап сборки
FROM golang:1.24-rc-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модуля и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник. CGO отключён для статической сборки.
RUN CGO_ENABLED=0 GOOS=linux go build -o backend

# Финальный образ
FROM scratch

# Копируем бинарник из сборочного этапа
COPY --from=builder /app/backend /backend

# Открываем порт приложения
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["/backend"]