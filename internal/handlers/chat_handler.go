package handlers

import (
	"net/http"

	"github.com/bytetiff/bytechat/internal/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Настраиваем WebSocket-соединение
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Обработчик WebSocket `/ws`
func ServeWs(hub *chat.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ⬇⬇⬇ Устанавливаем WebSocket-соединение ⬇⬇⬇
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
			return
		}

		// ⬇⬇⬇ Создаём нового клиента ⬇⬇⬇
		client := &chat.Client{
			Hub:  hub,
			Conn: conn,
			Send: make(chan []byte, 256), // Канал для отправки сообщений клиенту
		}

		// ⬇⬇⬇ Регистрируем клиента в Hub ⬇⬇⬇
		hub.Register <- client

		// Запускаем горутины для чтения и записи
		go client.WritePump()
		go client.ReadPump()
	}
}
