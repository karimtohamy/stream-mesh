package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"stream-mesh/streaming/internal/app"
	"stream-mesh/streaming/internal/bootstrap"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//create new app  config
	application, err := app.NewApp()
	if err != nil {
		fmt.Println("Error loading app")
		fmt.Println(err)
		log.Fatalf("failed to initialize app: %v", err)
	}
	// bootstrap and load all deps and services
	bootstrap.Init(ctx, application)

	//start newly created app
	if err := application.Start(ctx); err != nil {
		fmt.Println("Error starting app")
		fmt.Println(err)
		log.Fatalf("failed to initialize app: %v", err)

	}

	//signal app to shut down gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	application.ShutDown(ctx)
	fmt.Println("Shutting down")
}
