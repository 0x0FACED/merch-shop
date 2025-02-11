package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server stopped with err: %v", err)
	}
}
