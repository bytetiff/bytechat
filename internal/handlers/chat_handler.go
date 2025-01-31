package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bytetiff/bytechat/internal/models"
	"github.com/bytetiff/bytechat/internal/repository"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatRepo repository.ChatRepository
}

// NewChatHandler конструктор
func NewChatHandler(chatRepo repository.ChatRepository) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo}
}

type createChatRequest struct {
	Type    models.ChatType `json:"type" binding:"required"`
	Members []string        `json:"members"` // user_ids
}

// CreateChat - обработчик POST /chats
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req createChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Данные запроса:", req)

	// Создаем запись в chats
	chat := &models.Chat{
		Type: req.Type,
	}
	if err := h.chatRepo.CreateChat(context.Background(), chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Добавляем авторизованного пользователя в список участников
	// Добавляем текущего пользователя в список участников
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Удаляем дубликаты user_id
	uniqueMembers := make(map[string]bool)
	uniqueMembers[userID.(string)] = true // Гарантируем, что текущий юзер добавлен

	for _, member := range req.Members {
		uniqueMembers[member] = true
	}

	// Преобразуем в массив
	var members []string
	for k := range uniqueMembers {
		members = append(members, k)
	}

	// Добавляем пользователей в чат
	if err := h.chatRepo.AddMembers(context.Background(), chat.ID, members); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add members"})
		return
	}

	if err := h.chatRepo.AddMembers(context.Background(), chat.ID, members); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add members"})
		return
	}

	// Возвращаем созданный чат
	c.JSON(http.StatusOK, gin.H{
		"id":         chat.ID,
		"type":       chat.Type,
		"created_at": chat.CreatedAt,
		"members":    members,
	})
}

// GetUserChats - GET /chats
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chats, err := h.chatRepo.GetUserChats(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// AddMember - POST /chats/:chat_id/members
type addMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *ChatHandler) AddMember(c *gin.Context) {
	chatID := c.Param("chat_id")
	var req addMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.chatRepo.AddMembers(c.Request.Context(), chatID, []string{req.UserID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}
