package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ✅ Глобальная переменная DB
var DB *sql.DB

// InitDB - инициализация подключения к PostgreSQL
func InitDB() {
	dsn := "postgres://bytechat_user:secret@localhost:5432/bytechat?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к БД: %v", err))
	}

	if err = DB.Ping(); err != nil {
		panic(fmt.Sprintf("❌ Ошибка проверки соединения с БД: %v", err))
	}

	fmt.Println("✅ Подключено к PostgreSQL")
}
