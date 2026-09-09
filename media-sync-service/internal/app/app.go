package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stream-mesh/media-sync/internal/broker"
	"stream-mesh/media-sync/internal/config"
	"stream-mesh/media-sync/internal/storage"
	"stream-mesh/media-sync/internal/transmuxer"
)

type App struct {
	rabbitMQ   *broker.RabbitClient
	storage    storage.Storage
	transmuxer *transmuxer.Transmuxer
	cfg        *config.Config
}

func New(cfg *config.Config) (*App, error) {
	rabbitClient, err := broker.NewRabbitClient(cfg)
	if err != nil {
		fmt.Print("failed to start app due to rabbit mq")
		return nil, err
	}
	storageClient, err := storage.NewMinIOStorage(cfg)
	if err != nil {
		fmt.Print("failed to start app due to minio storage")
		return nil, err
	}
	return &App{
		rabbitMQ:   rabbitClient,
		storage:    storageClient,
		transmuxer: transmuxer.NewTransmuxer(),
		cfg:        cfg,
	}, nil
}

func (app *App) Start(ctx context.Context) error {

	rabbitTemplate := broker.NewRabbitTemplate(app.rabbitMQ)
	publisher := broker.NewPublisher(rabbitTemplate, app.cfg)
	listener := broker.NewListener(app.rabbitMQ, app.cfg)
	err := listener.StartTranscodeConsumer(ctx, func(ctx context.Context, job broker.TransCodeEvent) error {
		bucket := job.TargetBucket
		if bucket == "" {
			bucket = "vod-streams"
		}

		manifestKey := fmt.Sprintf("%s/master.m3u8", job.MediaId)
		manifestURL := fmt.Sprintf("/%s/%s/master.m3u8", bucket, job.MediaId)

		exists, err := app.storage.Exists(ctx, bucket, manifestKey)
		if err != nil {
			return fmt.Errorf("failed to check object existence: %w", err)
		}
		if exists {
			log.Printf("[DEDUP] Media %s already processed. Emitting completed event.", job.MediaId)
			return publisher.PublishTranscodeCompleted(ctx, job.MediaId, manifestURL, bucket)
		}

		workDir, err := os.MkdirTemp("", "transmux-"+job.MediaId+"-*")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer os.RemoveAll(workDir)

		sourceFile := filepath.Join(workDir, "source.mp4")
		if err := app.storage.Download(ctx, bucket, job.SourceFile, sourceFile); err != nil {
			return fmt.Errorf("failed to download source video %s: %w", job.SourceFile, err)
		}

		hlsDir := filepath.Join(workDir, "hls")
		if _, err := app.transmuxer.Transmux(ctx, sourceFile, hlsDir); err != nil {
			return fmt.Errorf("transmuxing failed: %w", err)
		}

		if err := app.storage.UploadDirectory(ctx, bucket, hlsDir, job.MediaId); err != nil {
			return fmt.Errorf("failed to upload HLS directory: %w", err)
		}

		log.Printf("[DONE] Successfully processed Media %s", job.MediaId)
		return publisher.PublishTranscodeCompleted(ctx, job.MediaId, manifestURL, bucket)
	})
	if err != nil {
		return err
	}
	defer app.rabbitMQ.Close()

	return nil
}
