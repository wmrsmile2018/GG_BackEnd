package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type ChatClients struct {
	user *model.User
	hub  *Hub
	conn *websocket.Conn
	send chan *model.Send
}

func NewClient(h *Hub, c *websocket.Conn, user *model.User, s chan *model.Send) *ChatClients {
	return &ChatClients{
		hub:  h,
		conn: c,
		user: user,
		send: s,
	}
}

func (c *ChatClients) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

	for {
		var message model.Message
		var send model.Send
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		json.Unmarshal(data, &message)
		send.User, err = c.hub.Store.User().Find(message.IdUser)
		send.Message = &message
		c.hub.Send <- &send

	}
}

func (c *ChatClients) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case send, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = json.Marshal(send)
			if err != nil {
				return
			}
			//w.Write(data)
			w.Write([]byte(send.Message.Message))
			//Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				_, err := json.Marshal(<-c.send)
				if err != nil {
					return
				}
				//w.Write(data)
				w.Write(send.Message.BytesMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
