package database

import (
 "database/sql"
 "fmt"

 _ "github.com/lib/pq"
 "apk-server/config"
)

var DB *sql.DB

// ConnectDB устанавливает подключение к БД, используя данные из конфигурации.
func ConnectDB(cfg *config.Config) (*sql.DB, error) {
 connStr := fmt.Sprintf(
  "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
  cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
 )
 db, err := sql.Open("postgres", connStr)
 if err != nil {
  return nil, err
 }
 if err = db.Ping(); err != nil {
  return nil, err
 }
 return db, nil
}

// SetDB сохраняет ссылку на подключение в глобальной переменной (при необходимости).
func SetDB(database *sql.DB) {
 DB = database
}