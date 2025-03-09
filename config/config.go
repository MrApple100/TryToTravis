package config

import (
	"os"
)

// Config содержит настройки приложения.
type Config struct {
	// Конфигурация базы данных.
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string

	// Конфигурация AWS S3/MinIO.
	AWSAccessKey string
	AWSSecretKey string
	AWSRegion    string
	S3Bucket     string
}

// LoadConfig загружает настройки из переменных окружения или использует значения по умолчанию.
func LoadConfig() *Config {
	return &Config{
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "secret"),
		DBName:       getEnv("DB_NAME", "apkdbZero"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		AWSAccessKey: getEnv("AWS_ACCESS_KEY", ""),
		AWSSecretKey: getEnv("AWS_SECRET_KEY", ""),
		AWSRegion:    getEnv("AWS_REGION", "us-east-1"),
		S3Bucket:     getEnv("S3_BUCKET", "apk-storage"),
	}
}

// getEnv возвращает значение переменной окружения key или fallback, если переменная не установлена.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
