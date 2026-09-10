package broker

import "time"

type TransCodeEvent struct {
	EventId      string    `json:"event_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	SourceFile   string    `json:"source_file"`
	TargetBucket string    `json:"target_bucket"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransmuxEventResponse struct {
	ManifestURL  string    `json:"manifest_url"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	ThumbnailURL string    `json:"thumbnail_url"`
	CompletedAt  time.Time `json:"completed_at"`
}
