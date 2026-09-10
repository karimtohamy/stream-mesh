package service

import (
	"context"
	"stream-mesh/streaming/internal/models"
	"stream-mesh/streaming/internal/repository"
)

type VideoService struct {
	repo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) GetVideoByID(id string) (*models.Video, error) {
	return s.repo.FindById(id)
}

func (s *VideoService) Save(ctx context.Context, video *models.Video) error {
	return s.repo.Save(ctx, video)
}

func (s *VideoService) GetBySlug(slug string) (*models.Video, error) {
	return s.repo.FindBySlug(slug)
}
