package main

import (
	"github.com/bytetiff/bytechat/internal/chat"
	"github.com/bytetiff/bytechat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ✅ Создаём Hub для управления WebSocket-клиентами
	hub := chat.NewHub()
	go hub.Run()

	// ✅ Регистрируем WebSocket-маршрут
	r.GET("/ws", handlers.ServeWs(hub))

	r.Run(":8080")
}
