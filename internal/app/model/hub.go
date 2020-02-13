package model

import "fmt"

type Hub struct {
	AllConnections 	map[*ChatClients]bool
	UserConnection	map[string]*ChatClients
	Message			chan []byte
	Register		chan *ChatClients
	Unregister		chan *ChatClients
	Clients			chan *User
	IsGeneralChat	bool
}

func NewHub() *Hub {
	return &Hub{
		Message:        make(chan []byte), // broadcast
		Register:       make(chan *ChatClients),
		Unregister:     make(chan *ChatClients),
		AllConnections: make(map[*ChatClients]bool), // clients
		UserConnection: make(map[string]*ChatClients), // clients
		Clients:		make(chan *User),
		IsGeneralChat:  false,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.AllConnections[client] = true
		case client := <-h.Unregister:
			if _, ok := h.AllConnections[client]; ok {
				delete(h.AllConnections, client)
				close(client.send)
			}
		case message := <-h.Message:
			if h.IsGeneralChat {
				for conn := range h.AllConnections {
					fmt.Println(conn, message)
					select {
					case conn.send <- message:
					default:
						close(conn.send)
						delete(h.AllConnections, conn)
					}
				}
			} else {
				for client := range h.Clients {
					conn := h.UserConnection[client.ID]
					fmt.Println(conn, client.ID)
					select {
					case conn.send <- message:
					default:
						close(conn.send)
						delete(h.UserConnection, client.ID)
					}
				}
			}

		}
	}
}
