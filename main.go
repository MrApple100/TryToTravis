package main

import (
 "log"

 "apk-server/config"
 "apk-server/database"
 "apk-server/handlers"
 "apk-server/storage"

 "github.com/gin-gonic/gin"
)

func main() {
 // Загружаем конфигурацию из переменных окружения.
 cfg := config.LoadConfig()

 // Подключаемся к базе данных с использованием конфигурации.
 db, err := database.ConnectDB(cfg)
 if err != nil {
  log.Fatalf("Не удалось подключиться к БД: %v", err)
 }
 database.SetDB(db)

 // Инициализируем клиента S3.
 storage.InitS3(cfg)

 // Инициализируем HTTP-сервер с использованием Gin.
 r := gin.Default()

 api := r.Group("/api/apk")
 {
  api.POST("/upload", handlers.UploadApk)
  api.GET("/latest", handlers.GetLatestVersion)
  api.GET("/download", handlers.DownloadApk)
 }

 // Запускаем сервер на порту 8080.
 r.Run(":8080")
}