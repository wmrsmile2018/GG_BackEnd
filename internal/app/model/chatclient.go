package model

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
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

type Chat struct {
	TypeChat		string
	BytesMessage	[]byte
	Message 		string
}
type ChatClients struct {
	user 	*User
	hub 	*Hub
	conn 	*websocket.Conn
	send 	chan []byte
}



func NewClient(h *Hub, c *websocket.Conn, s chan[]byte, user *User) *ChatClients {
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
		var chat Chat
		_, message, err := c.Connection.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		json.Unmarshal(message, &chat)
		chat.BytesMessage = []byte(chat.Message)
		chat.BytesMessage = bytes.TrimSpace(bytes.Replace(chat.BytesMessage, newline, space, -1))
		//fmt.Println(chat.BytesMessage)
		//fmt.Println(string(chat.BytesMessage))
		//fmt.Println(chat.Message)
		//fmt.Println(chat.TypeChat)
		c.Connection.hub.Message <- chat.BytesMessage
		//go func(mes []byte) {c.hub.Message <- mes}(chat.BytesMessage)
		//c.hub.Message <- chat.BytesMessage
		c.Connection.hub.IsGeneralChat = chat.TypeChat == "general"
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
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Connection.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Connection.send)
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
