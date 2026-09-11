package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Slug        string `gorm:"uniqueIndex;not null;size:255"`
	ManifestURL string `gorm:"not null;size:255"`
	Thumbnail   string `gorm:"not null;size:255"`
}
