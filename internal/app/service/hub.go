package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/store"
	"time"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Hub struct {
	AllConnections map[*ChatClients]bool
	Send           chan *model.Send
	Register       chan *ChatClients
	Unregister     chan *ChatClients
	Store          store.Store
	AllConns       map[string]map[*ChatClients]bool
}

func NewHub(store store.Store) *Hub {
	return &Hub{
		Register:       make(chan *ChatClients),
		Unregister:     make(chan *ChatClients),
		AllConnections: make(map[*ChatClients]bool), // clients
		Store:          store,
		AllConns:       make(map[string]map[*ChatClients]bool),
		Send:           make(chan *model.Send),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.AllConnections[client] = true
			fmt.Println("connections:  ", len(h.AllConnections))
		case client := <-h.Unregister:
			if _, ok := h.AllConnections[client]; ok {
				delete(h.AllConns[client.user.ID], client)
				delete(h.AllConnections, client)
				close(client.send)
			}
		case send := <-h.Send:
			send.Message.IdMessage = uuid.New().String()
			send.Message.BytesMessage = []byte(send.Message.Message)
			send.Message.BytesMessage = bytes.TrimSpace(bytes.Replace(send.Message.BytesMessage, newline, space, -1))
			send.Message.TimeCreateM = time.Now().Unix()

			if send.Message.TypeChat == "general" {
				//h.Store.User().CreateMessage(send.Message)
				fmt.Println(h.Store.User().CreateMessage(send.Message))
				if len(h.AllConnections) != 0 {
					for conn := range h.AllConnections {
						var s model.Send
						s.Message = send.Message
						s.User = conn.user
						select {
						case conn.send <- &s:
						default:
							close(conn.send)
							delete(h.AllConns[conn.user.ID], conn)
							delete(h.AllConnections, conn)
						}
					}
				}
			} else {
				users := getUsers(h, send.Message.IdChat)
				if users != nil && users[send.User.ID] {
					fmt.Println("i am here")
					fmt.Println(h.Store.User().CreateMessage(send.Message))
					for user := range users {
						connections := h.AllConns[user]
						if len(connections) != 0 {
							for conn := range connections {
								fmt.Println(conn, conn.user)
								var s model.Send
								s.Message = send.Message
								s.User = conn.user
								select {
								case conn.send <- &s:
								default:
									close(conn.send)
									delete(h.AllConns[conn.user.ID], conn)
									delete(h.AllConnections, conn)
								}
							}
						}
					}
				}
			}
		}
	}
}

func getUsers(h *Hub, id string) map[string]bool {
	users, err := h.Store.User().FindByChat(id)
	if err != nil {
		return nil
	}
	return users
}
