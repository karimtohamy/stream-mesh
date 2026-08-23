package broker

import "time"

type TransCodeEvent struct {
	Eventid      string    `json:"event_id"`
	MediaId      string    `json:"media_id"`
	SourceFile   string    `json:"source_file"`
	TargetBucket string    `json:"target_bucket"`
	Resolution   string    `json:"resolution"`
	Checksum     string    `json:"checksum"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransmuxEventResponse struct {
	MediaId string `json:"media_id"`
	ManifestURL string `json:"manifest_url"`
	TargetBucket string `json:"target_bucket"`
	CompletedAt time.Time `json:"completed_at"`
}

type UsageTickPayload struct {
	EventID         string    `json:"eventId"`
	StreamID        string    `json:"streamId"`
	UserID          string    `json:"userId"`
	DurationSeconds int       `json:"durationSeconds"`
	CostPerMinute   float64   `json:"costPerMinute"`
	Timestamp       time.Time `json:"timestamp"`
}