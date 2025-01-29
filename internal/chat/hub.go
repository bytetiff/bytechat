package chat

// Hub управляет WebSocket-клиентами
type Hub struct {
	Clients    map[*Client]bool // Список всех клиентов
	Broadcast  chan []byte      // Канал для рассылки сообщений
	Register   chan *Client     // Клиенты, которые подключились
	Unregister chan *Client     // Клиенты, которые отключились
}

// Создаём новый Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Запускаем Hub в отдельной горутине
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
