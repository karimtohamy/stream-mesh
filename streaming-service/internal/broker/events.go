package broker

import "time"

type NewVideoEvent struct {
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	ManifestURL  string    `json:"manifest_url"`
	ThumbNailURL string    `json:"thumbnail_url"`
	CompletedAt  time.Time `json:"completed_at"`
}
