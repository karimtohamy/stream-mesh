package ws

import "time"

type Message struct {
	RoomCode  string    `json:"room_code"`
	UserId    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Data      any       `json:"data"`
}

