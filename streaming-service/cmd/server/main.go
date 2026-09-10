package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	app2 "stream-mesh/streaming/internal/app"
	"stream-mesh/streaming/internal/config"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := app2.NewApp(cfg)
	if err != nil {
		fmt.Println("Error loading app")
		fmt.Println(err)
		log.Fatalf("failed to initialize app: %v", err)
	}
	if err := app.Start(ctx); err != nil {
		fmt.Println("Error starting app")
		fmt.Println(err)

		log.Fatalf("failed to initialize app: %v", err)

	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	app.ShutDown(ctx)
	fmt.Println("Shutting down")
}
