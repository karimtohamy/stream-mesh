package service

import (
	"context"
	"fmt"
	"log"
	"stream-mesh/streaming/internal/cache"
	"stream-mesh/streaming/internal/models"
	"stream-mesh/streaming/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type VideoService struct {
	repo  repository.VideoRepository
	redis *redis.Client
}

func NewVideoService(repo repository.VideoRepository, rdb *redis.Client) *VideoService {
	return &VideoService{
		repo:  repo,
		redis: rdb,
	}
}

func (s *VideoService) GetVideoByID(ctx *gin.Context, id string) (*models.Video, error) {
	cacheKey := fmt.Sprintf("video:id:%s", id)
	if cached, _ := cache.Get[models.Video](ctx, s.redis, cacheKey); cached != nil {
		log.Printf("CACHE HIT %s", cacheKey)
		return cached, nil
	}
	log.Printf("CACHE MISS FOR KEY %s", cacheKey)
	video, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(ctx, s.redis, cacheKey, video, time.Hour); err != nil {
		log.Printf("CACHE SET FAIL FOR KEY %s", cacheKey)
		return nil, err
	}
	log.Printf("CACHE SET %s", cacheKey)
	return video, nil
}
func (s *VideoService) Save(ctx context.Context, video *models.Video) error {
	if err := s.repo.Save(ctx, video); err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("video:id:%d", video.ID)
	if err := cache.Set(ctx, s.redis, cacheKey, *video, time.Hour); err != nil {
		log.Printf("[CACHE ERROR] failed to set %s: %v", cacheKey, err)
	}

	log.Printf("[CACHE SET] %s", cacheKey)
	return nil
}

func (s *VideoService) GetBySlug(ctx *gin.Context, slug string) (*models.Video, error) {
	cacheKey := fmt.Sprintf("video:slug:%s", slug)
	if cached, _ := cache.Get[models.Video](ctx, s.redis, cacheKey); cached != nil {
		log.Printf("CACHE HIT %s", cacheKey)
		return cached, nil
	}
	log.Printf("CACHE MISS FOR KEY %s", cacheKey)
	video, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(ctx, s.redis, cacheKey, video, time.Hour); err != nil {
		log.Printf("CACHE SET FAIL FOR KEY %s", cacheKey)
		return nil, err
	}
	log.Printf("CACHE SET %s", cacheKey)
	return video, nil
}
