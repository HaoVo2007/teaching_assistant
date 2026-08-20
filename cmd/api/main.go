package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"teaching_assistant/internal/app"
	"teaching_assistant/internal/config"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	application, err := app.NewApplication(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = application.Shutdown(ctx)
}
