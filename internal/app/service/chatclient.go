package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"log"
	"time"
)

const (
	writeWait 		= 10 * time.Second
	pongWait 		= 60 * time.Second
	pingPeriod 		= (pongWait * 9) / 10
	maxMessageSize 	= 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type ChatClients struct {
	user 	*model.User
	hub 	*Hub
	conn 	*websocket.Conn
	send 	chan *model.Message
}



func NewClient(h *Hub, c *websocket.Conn, s chan *model.Message, user *model.User) *ChatClients {
	return &ChatClients{
		hub:  h,
		conn: c,
		send: s,
		user: user,
	}
}

func (c *ClientConn) ReadPump() {
	defer func() {
		c.Connection.hub.Unregister <- c
		c.Connection.conn.Close()
	}()
	c.Connection.conn.SetReadLimit(maxMessageSize)
	c.Connection.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.conn.SetPongHandler(func(string) error { c.Connection.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var mes model.Message
		_, message, err := c.Connection.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		json.Unmarshal(message, &mes)
		users , err := c.Connection.hub.Store.User().FindByChat(mes.IdChat)
		if err != nil {
			return
		}
		for client := range c.Connection.hub.Users {
			fmt.Println("__________________client ",client)
		}
		mes.BytesMessage = []byte(mes.Message)
		mes.BytesMessage = bytes.TrimSpace(bytes.Replace(mes.BytesMessage, newline, space, -1))
		mes.User, err =	c.Connection.hub.Store.User().Find(mes.IdUser)
		if err != nil {
			return
		}
		mes.TimeCreateM = time.Now().UnixNano()
		mes.IdMessage = uuid.New().String()
		c.Connection.hub.Users = users
		c.Connection.hub.Message <- &mes
		c.Connection.hub.IsGeneralChat = (mes.TypeChat == "general")
	}
}

func (c *ClientConn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Connection.send:
			c.Connection.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Connection.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Connection.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
 			data, err := json.Marshal(message)
 			if err != nil {
				return
			}
			_, err = c.Connection.hub.Store.User().CreateMessage(message)
			if err != nil {
				return
			}
			//w.Write(data)
			fmt.Println(data)
			w.Write(message.BytesMessage)
			// Add queued chat messages to the current websocket message.
			n := len(c.Connection.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				data, err := json.Marshal(<-c.Connection.send)
				if err != nil {
					return
				}
				//w.Write(data)
				fmt.Println(data)
				w.Write(message.BytesMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Connection.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
