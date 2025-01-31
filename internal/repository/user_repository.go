package repository

import (
	"context"

	"github.com/bytetiff/bytechat/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error) // ДОЛЖЕН БЫТЬ!
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db}
}

// Create сохраняет пользователя в БД
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `
	row := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash)
	if err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}
	return nil
}

// GetByID - получение пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
        SELECT id, username, email, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1
    `
	row := r.db.QueryRow(ctx, query, id)

	var user models.User
	if err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail ищет пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, username, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	row := r.db.QueryRow(ctx, query, email)

	var user models.User
	if err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
