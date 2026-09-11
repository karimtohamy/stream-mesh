package ws

import "encoding/json"

type Hub struct {
	rooms      map[string]map[*Client]bool // roomCode -> set of connected clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// create the room bucket on first join
			if h.rooms[client.roomCode] == nil {
				h.rooms[client.roomCode] = make(map[*Client]bool)
			}
			h.rooms[client.roomCode][client] = true

		case client := <-h.unregister:
			if room, ok := h.rooms[client.roomCode]; ok {
				delete(room, client)
				close(client.send)
				// clean up empty rooms so we don't leak memory
				if len(room) == 0 {
					delete(h.rooms, client.roomCode)
				}
			}

		case msg := <-h.broadcast:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			for client := range h.rooms[msg.RoomCode] {
				select {
				case client.send <- data:
				default:
					// client send buffer is full — connection is dead, drop it
					close(client.send)
					delete(h.rooms[msg.RoomCode], client)
				}
			}
		}
	}
}
func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- msg
}

