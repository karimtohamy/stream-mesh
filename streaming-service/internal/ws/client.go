package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	roomCode string
	userId   string
}

func NewClient(hub *Hub, conn *websocket.Conn, roomCode, userId string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		roomCode: roomCode,
		userId:   userId,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		if err := c.conn.Close(); err != nil {
			log.Printf("[WS] close error for user %s: %v", c.userId, err)
		}
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		msg.UserId = c.userId
		msg.RoomCode = c.roomCode
		c.hub.Broadcast(&msg)
	}
}

func (c *Client) WritePump() {
	defer func() {
		if err := c.conn.Close(); err != nil {
			log.Printf("[WS] close error for user %s: %v", c.userId, err)
		}
	}()
	for data := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}
