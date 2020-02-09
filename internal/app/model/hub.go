package model

type Hub struct {
	Clients map[*ChatClients]bool
	Broadcast chan []byte
	Register chan *ChatClients
	Unregister chan *ChatClients
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *ChatClients),
		Unregister: make(chan *ChatClients),
		Clients:    make(map[*ChatClients]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
