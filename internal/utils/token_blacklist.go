package utils

import (
	"fmt"
	"sync"
	"time"
)

// ✅ Теперь мы используем глобальную `DB` из db.go
var mu sync.Mutex

// AddToBlacklist - добавляет токен в "чёрный список" (PostgreSQL)
func AddToBlacklist(token string) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("🚫 Добавляем токен в чёрный список (PostgreSQL):", token)

	query := `
        INSERT INTO token_blacklist (token, created_at)
        VALUES ($1, $2)
        ON CONFLICT (token) DO NOTHING
    `
	_, err := DB.Exec(query, token, time.Now())
	if err != nil {
		fmt.Println("❌ Ошибка добавления токена в blacklist:", err)
	}
}

// IsBlacklisted - проверяет, заблокирован ли токен (PostgreSQL)
func IsBlacklisted(token string) bool {
	mu.Lock()
	defer mu.Unlock()

	query := `SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE token = $1)`
	var exists bool
	err := DB.QueryRow(query, token).Scan(&exists)
	if err != nil {
		fmt.Println("❌ Ошибка проверки blacklist:", err)
		return false
	}

	fmt.Printf("🔍 Проверка токена '%s' в blacklist → %v\n", token, exists)
	return exists
}
