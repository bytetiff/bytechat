package repositories

import (
	"database/sql"
	"fmt"

	"github.com/bytetiff/bytechat/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser - добавить пользователя в БД
func (r *UserRepository) CreateUser(u models.User) (models.User, error) {
	query := `INSERT INTO users (username, password, created_at)
              VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(query, u.Username, u.Password, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("❌ Ошибка при вставке пользователя: %v", err)
	}
	return u, nil
}

// GetUserByUsername - получить пользователя по логину
func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, created_at
	          FROM users
	          WHERE username = $1 LIMIT 1`

	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
