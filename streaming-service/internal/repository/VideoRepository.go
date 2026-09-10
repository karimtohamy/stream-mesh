package repository

import (
	"context"
	"stream-mesh/streaming/internal/models"

	"gorm.io/gorm"
)

type videoRepository struct {
	db *gorm.DB
}
type VideoRepository interface {
	Save(ctx context.Context, video *models.Video) error
	FindBySlug(slug string) (*models.Video, error)
	FindById(ID string) (*models.Video, error)
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db}
}
func (r *videoRepository) Save(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) FindBySlug(slug string) (*models.Video, error) {
	var video models.Video
	err := r.db.Where("slug = ?", slug).First(&video).Error
	return &video, err
}

func (r *videoRepository) FindById(id string) (*models.Video, error) {
	var video models.Video
	err := r.db.Where("id = ?", id).First(&video).Error
	return &video, err
}
