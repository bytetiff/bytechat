package repository

import (
	"context"
	"errors"

	"github.com/bytetiff/bytechat/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	AddMembers(ctx context.Context, chatID string, userIDs []string) error
	GetUserChats(ctx context.Context, userID string) ([]models.Chat, error)
	IsMember(ctx context.Context, chatID, userID string) (bool, error)
}

type chatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepository {
	return &chatRepository{db}
}

// CreateChat создает новую запись в chats
func (r *chatRepository) CreateChat(ctx context.Context, chat *models.Chat) error {
	query := `
        INSERT INTO chats (type)
        VALUES ($1)
        RETURNING id, created_at
    `
	err := r.db.QueryRow(ctx, query, chat.Type).Scan(&chat.ID, &chat.CreatedAt)
	return err
}

// AddMembers добавляет участников в chat_members
func (r *chatRepository) AddMembers(ctx context.Context, chatID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return errors.New("no user ids provided")
	}

	// Подготавливаем батч
	batch := &pgx.Batch{}
	for _, userID := range userIDs {
		batch.Queue(`
            INSERT INTO chat_members (chat_id, user_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING
        `, chatID, userID)
	}

	br := r.db.SendBatch(ctx, batch)
	_, err := br.Exec()
	_ = br.Close()
	return err
}

// GetUserChats возвращает список чатов, в которых участвует пользователь
func (r *chatRepository) GetUserChats(ctx context.Context, userID string) ([]models.Chat, error) {
	query := `
        SELECT c.id, c.type, c.created_at
        FROM chats c
        INNER JOIN chat_members cm ON c.id = cm.chat_id
        WHERE cm.user_id = $1
        ORDER BY c.created_at DESC
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}

func (r *chatRepository) IsMember(ctx context.Context, chatID, userID string) (bool, error) {
	query := `
        SELECT COUNT(*) FROM chat_members
        WHERE chat_id = $1 AND user_id = $2
    `
	var count int
	err := r.db.QueryRow(ctx, query, chatID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
