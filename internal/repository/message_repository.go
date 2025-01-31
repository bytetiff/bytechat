package repository

import (
	"context"

	"github.com/bytetiff/bytechat/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *models.Message) error
	GetMessagesByChat(ctx context.Context, chatID string, limit, offset int) ([]models.Message, error)
}

type messageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) MessageRepository {
	return &messageRepository{db: db}
}

// CreateMessage - сохраняет сообщение в БД
func (r *messageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	query := `
        INSERT INTO messages (chat_id, sender_id, content)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	err := r.db.QueryRow(ctx, query, msg.ChatID, msg.SenderID, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
	return err
}

// GetMessagesByChat - возвращает список сообщений в чате с пагинацией
func (r *messageRepository) GetMessagesByChat(ctx context.Context, chatID string, limit, offset int) ([]models.Message, error) {
	query := `
        SELECT id, chat_id, sender_id, content, created_at
        FROM messages
        WHERE chat_id = $1
        ORDER BY created_at ASC
        LIMIT $2 OFFSET $3
    `
	rows, err := r.db.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}
