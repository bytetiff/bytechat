package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshTokenRepository управляет refresh-токенами в БД
type RefreshTokenRepository interface {
	SaveToken(ctx context.Context, userID string, token string, expiresAt time.Time) error
	DeleteToken(ctx context.Context, token string) error
	GetUserByToken(ctx context.Context, token string) (string, error)
}

type refreshTokenRepository struct {
	db *pgxpool.Pool
}

// NewRefreshTokenRepository - конструктор
func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

// SaveToken - сохраняет refresh-токен в БД
func (r *refreshTokenRepository) SaveToken(ctx context.Context, userID string, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, token, expiresAt)
	return err
}

// DeleteToken - удаляет refresh-токен
func (r *refreshTokenRepository) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(ctx, query, token)
	return err
}

// GetUserByToken - ищет пользователя по refresh-токену
func (r *refreshTokenRepository) GetUserByToken(ctx context.Context, token string) (string, error) {
	query := `SELECT user_id FROM refresh_tokens WHERE token = $1`
	var userID string
	err := r.db.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
