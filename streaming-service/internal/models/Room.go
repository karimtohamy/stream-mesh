package models

import "time"

type Room struct {
	Code          string         `json:"code"`
	Members       []string       `json:"members"`
	CurrentVideo  *string        `json:"current_video,omitempty"`
	PlaybackState *PlaybackState `json:"playback_state,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
type PlaybackState struct {
	Position float64 `json:"position"`
	Playing  bool    `json:"playing"`
}
