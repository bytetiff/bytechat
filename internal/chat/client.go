package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Клиент WebSocket-соединения
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// Читаем из WebSocket и отправляем в `Hub.Broadcast`

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg ChatMessage
		err = json.Unmarshal(rawMessage, &msg)
		if err != nil {
			fmt.Println("❌ Ошибка парсинга JSON:", err)
			continue
		}

		formattedMessage, _ := json.Marshal(msg) // Обратно в JSON
		c.Hub.Broadcast <- formattedMessage
	}
}

// Читаем из `Send` и отправляем в WebSocket клиенту
func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
