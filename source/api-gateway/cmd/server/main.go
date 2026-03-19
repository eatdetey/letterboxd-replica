package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("init gateway: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		os.Exit(1)
	}
}
