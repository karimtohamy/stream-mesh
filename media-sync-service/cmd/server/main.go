package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"stream-mesh/media-sync/internal/broker"
	"stream-mesh/media-sync/internal/config"
	"stream-mesh/media-sync/internal/storage"
	"stream-mesh/media-sync/internal/transmuxer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rabbitClient, err := broker.NewRabbitClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	storageClient, err := storage.NewMinIOStorage(&cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	rabbitTemplate := broker.NewRabbitTemplate(rabbitClient)
	publisher := broker.NewPublisher(rabbitTemplate, cfg)
	listener := broker.NewListener(rabbitClient, cfg)
	engine := transmuxer.NewTransmuxer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = listener.StartTranscodeConsumer(ctx, func(ctx context.Context, job broker.TransCodeEvent) error {
		bucket := job.TargetBucket
		if bucket == "" {
			bucket = "vod-streams"
		}

		manifestKey := fmt.Sprintf("%s/master.m3u8", job.MediaId)
		manifestURL := fmt.Sprintf("/%s/%s/master.m3u8", bucket, job.MediaId)

		exists, err := storageClient.Exists(ctx, bucket, manifestKey)
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
		if err := storageClient.Download(ctx, bucket, job.SourceFile, sourceFile); err != nil {
			return fmt.Errorf("failed to download source video %s: %w", job.SourceFile, err)
		}

		hlsDir := filepath.Join(workDir, "hls")
		if _, err := engine.Transmux(ctx, sourceFile, hlsDir); err != nil {
			return fmt.Errorf("transmuxing failed: %w", err)
		}

		if err := storageClient.UploadDirectory(ctx, bucket, hlsDir, job.MediaId); err != nil {
			return fmt.Errorf("failed to upload HLS directory: %w", err)
		}

		log.Printf("[DONE] Successfully processed Media %s", job.MediaId)
		return publisher.PublishTranscodeCompleted(ctx, job.MediaId, manifestURL, bucket)
	})
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Printf("Worker online. Listening on %s...", cfg.RabbitMQ.TranscodeQueue)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker gracefully...")
}
