package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword - хеширует пароль перед сохранением в БД
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword - сравнивает хеш пароля с введённым паролем
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
