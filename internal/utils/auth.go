package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи токенов
var secretKey = []byte("supersecretkey")

// GenerateToken - создаёт JWT-токен для пользователя
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 часа жизни
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken - проверка токена (проверяем чёрный список + подпись)
func ParseToken(tokenString string) (string, error) {
	// Проверяем, не заблокирован ли токен
	if IsBlacklisted(tokenString) {
		return "", errors.New("token has been revoked")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, exists := claims["username"].(string)
		if !exists {
			return "", errors.New("invalid token data")
		}
		return username, nil
	}

	return "", errors.New("invalid token")
}

// ExtractToken - извлекает токен из заголовка Authorization (Bearer ...)
func ExtractToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header")
	}
	return parts[1], nil
}
