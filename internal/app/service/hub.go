package service

import (
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/store"
)

type Hub struct {
	AllConnections map[*ClientConn]bool
	UserConnection map[string]*ClientConn
	Message        chan *model.Message
	Register       chan *ClientConn
	Unregister     chan *ClientConn
	Users          map[*model.User]bool
	IsGeneralChat  bool
	text           string
	Store          store.Store
}


type ClientConn struct {
	Id			string
	Connection	*ChatClients
}

func NewHub(store store.Store) *Hub {
	return &Hub{
		Message:        make(chan *model.Message),
		Register:       make(chan *ClientConn),
		Unregister:     make(chan *ClientConn),
		AllConnections: make(map[*ClientConn]bool), // clients
		UserConnection: make(map[string]*ClientConn),
		Users:          make(map[*model.User]bool),
		IsGeneralChat:  false,
		Store:          store,
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
				close(client.Connection.send)
			}
		case message := <-h.Message:
			if h.IsGeneralChat {
				for conn := range h.AllConnections {
					select {
					case conn.Connection.send <- message:
					default:
						close(conn.Connection.send)
						delete(h.AllConnections, conn)
					}
				}
			} else {
				for user := range h.Users {
					conn := h.UserConnection[user.ID]
					if conn == nil {
						continue
					} else {
						select {
						case conn.Connection.send <-message:
						default:
							close(conn.Connection.send)
							delete(h.AllConnections, conn)
						}
					}
				}
			}
		}
	}
}
