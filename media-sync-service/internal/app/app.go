package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stream-mesh/media-sync/internal/broker"
	"stream-mesh/media-sync/internal/config"
	"stream-mesh/media-sync/internal/engine"
	"stream-mesh/media-sync/internal/storage"
)

type App struct {
	rabbitMQ *broker.RabbitClient
	storage  storage.Storage
	cfg      *config.Config
}

func New(cfg *config.Config) (*App, error) {
	rabbitClient, err := broker.NewRabbitClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	storageClient, err := storage.NewMinIOStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}
	return &App{
		rabbitMQ: rabbitClient,
		storage:  storageClient,
		cfg:      cfg,
	}, nil
}

func (a *App) Start(ctx context.Context) error {

	rabbitTemplate := broker.NewRabbitTemplate(a.rabbitMQ)
	publisher := broker.NewPublisher(rabbitTemplate, a.cfg)
	listener := broker.NewListener(a.rabbitMQ, a.cfg)

	return listener.StartTranscodeConsumer(ctx, func(ctx context.Context, job broker.TransCodeEvent) error {
		return a.processTranscodeJob(ctx, job, publisher)
	})
}

func (a *App) processTranscodeJob(ctx context.Context, job broker.TransCodeEvent, publisher *broker.Publisher) error {
	bucket := job.TargetBucket
	if bucket == "" {
		bucket = "vod-streams"
	}

	manifestKey := fmt.Sprintf("%s/hls/master.m3u8", job.Slug)
	manifestURL := fmt.Sprintf("/%s/%s/hls/master.m3u8", bucket, job.Slug)

	exists, err := a.storage.Exists(ctx, bucket, manifestKey)
	if err != nil {
		return fmt.Errorf("failed to check object existence: %w", err)
	}
	if exists {
		log.Printf("[DEDUP] %s already processed", job.Slug)
		return publisher.PublishTranscodeCompleted(ctx, manifestURL, "", job)
	}

	workDir, err := os.MkdirTemp("", "transmux-"+job.Slug+"-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	sourceFile := filepath.Join(workDir, "source.mp4")
	if err := a.storage.Download(ctx, bucket, job.SourceFile, sourceFile); err != nil {
		return fmt.Errorf("failed to download source: %w", err)
	}

	hlsDir := filepath.Join(workDir, "hls")
	if _, err := engine.Transmux(ctx, sourceFile, hlsDir); err != nil {
		return fmt.Errorf("transmuxing failed: %w", err)
	}

	thumbnailDir := filepath.Join(workDir, "thumbnail.jpg")
	if err := engine.GenerateThumbnail(ctx, sourceFile, thumbnailDir); err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	if err := a.storage.UploadDirectory(ctx, bucket, hlsDir, fmt.Sprintf("%s/hls", job.Slug)); err != nil {
		return fmt.Errorf("failed to upload HLS: %w", err)
	}
	if err := a.storage.Upload(ctx, bucket, fmt.Sprintf("%s/thumbnail/thumbnail.jpg", job.Slug), thumbnailDir); err != nil {
		return fmt.Errorf("failed to upload thumbnail: %w", err)
	}
	thumbnailURL := fmt.Sprintf("/%s/%s/thumbnail/thumbnail.jpg", bucket, job.Slug)
	log.Printf("[DONE] %s processed successfully", job.Slug)
	return publisher.PublishTranscodeCompleted(ctx, manifestURL, thumbnailURL, job)
}
func (a *App) Shutdown() {
	a.rabbitMQ.Close()
}
