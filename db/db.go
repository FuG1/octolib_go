package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

var DB *sql.DB

func InitDB() error {
	// Настройка строки подключения
	connStr := "host=# port=# user=# password=# dbname=# sslmode=disable"

	// Открытие соединения
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка соединения
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connected successfully")
	return nil
}
