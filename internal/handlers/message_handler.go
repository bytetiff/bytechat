package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bytetiff/bytechat/internal/models"
	"github.com/bytetiff/bytechat/internal/repository"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	msgRepo  repository.MessageRepository
	chatRepo repository.ChatRepository
}

func NewMessageHandler(msgRepo repository.MessageRepository, chatRepo repository.ChatRepository) *MessageHandler {
	return &MessageHandler{
		msgRepo:  msgRepo,
		chatRepo: chatRepo,
	}
}

// POST /chats/:chat_id/messages
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	chatID := c.Param("chat_id")

	// Приведение user_id к строке
	userIDRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
		return
	}

	// Проверяем, состоит ли userID в чате
	isMember, err := h.isChatMember(c.Request.Context(), chatID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this chat"})
		return
	}

	// Читаем JSON
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаём сообщение
	msg := &models.Message{
		ChatID:   chatID,
		SenderID: userID,
		Content:  req.Content,
	}
	if err := h.msgRepo.CreateMessage(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// GET /chats/:chat_id/messages?limit=20&offset=0
func (h *MessageHandler) GetMessages(c *gin.Context) {
	chatID := c.Param("chat_id")

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Проверяем, состоит ли userID в чате
	isMember, err := h.isChatMember(c.Request.Context(), chatID, userID.(string))
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this chat"})
		return
	}

	// Параметры limit и offset
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")

	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	messages, err := h.msgRepo.GetMessagesByChat(c.Request.Context(), chatID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// isChatMember - вспомогательная функция (проверяет, есть ли user в chat_members)
func (h *MessageHandler) isChatMember(ctx context.Context, chatID, userID string) (bool, error) {
	isMember, err := h.chatRepo.IsMember(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	return isMember, nil
}
