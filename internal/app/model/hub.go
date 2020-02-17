package model

import "fmt"

type Hub struct {
	AllConnections 	map[*ClientConn]bool
	UserConnection 	map[string]*ClientConn
	//ConnectionUser	map[*ClientConn]*User
	Message			chan []byte
	Register		chan *ClientConn
	Unregister		chan *ClientConn
	Clients			map[*User]bool
	IsGeneralChat	bool
}

type ClientConn struct {
	Id			string
	Connection	*ChatClients
}

func NewHub() *Hub {
	return &Hub{
		Message:        	make(chan []byte), // broadcast
		Register:       	make(chan *ClientConn),
		Unregister:     	make(chan *ClientConn),
		AllConnections: 	make(map[*ClientConn]bool), // clients
		UserConnection: 	make(map[string]*ClientConn),
		//ConnectionUser:		make(map[*ClientConn]*User),
		Clients:			make(map[*User]bool),
		IsGeneralChat:  	false,
	}
}

func (h *Hub) Run() {
	for {
		//fmt.Println("wtf___________________________________1")
		select {
		case client := <-h.Register:
			h.AllConnections[client] = true
			//fmt.Println("wtf___________________________________2")
		case client := <-h.Unregister:
			if _, ok := h.AllConnections[client]; ok {
				delete(h.AllConnections, client)
				close(client.Connection.send)
			}
			//fmt.Println("wtf___________________________________3")
		case message := <-h.Message:
			//fmt.Println("wtf___________________________________4")
			if h.IsGeneralChat {
				for conn := range h.AllConnections {
					//fmt.Println(conn)
					select {
					case conn.Connection.send <- message:
					default:
						close(conn.Connection.send)
						delete(h.AllConnections, conn)
					}
				}
			} else {
				for user := range h.Clients {
					conn := h.UserConnection[user.ID]
					if conn == nil {
						continue
					} else {
						fmt.Println("hub____________user", user)
						fmt.Println("hub____________conn", conn)
						fmt.Println("hub____________message", message)
						select {
						case conn.Connection.send <-message:
						default:
							close(conn.Connection.send)
							delete(h.AllConnections, conn)
						}
					}
				}
			}
			//fmt.Println("FINISH____________________________________")
		}
	}
}
