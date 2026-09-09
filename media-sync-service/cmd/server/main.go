package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"stream-mesh/media-sync/internal/app"
	"syscall"

	"stream-mesh/media-sync/internal/config"
)

func main() {
	//load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	//load context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//initialize app with config and connection
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	//start app
	err = application.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}

	//signal to shut down gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker gracefully...")
}
