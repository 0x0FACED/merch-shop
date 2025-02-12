package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/server"
	"github.com/0x0FACED/merch-shop/internal/server/handler"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.New(cfg.Logger)

	db, err := database.New(cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to create database", zap.Error(err))
	}
	db.MustConnect(ctx)
	defer db.Close()

	merchService := service.NewUserService(db, log)

	h := handler.NewHandler(merchService, log, &cfg.Server)

	server, err := server.NewServer(cfg, h)
	if err != nil {
		log.Fatal("Failed to create server", zap.Error(err))
	}

	go func() {
		if err := server.Start(ctx); err != nil {
			log.Error("Server stopped with error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down gracefully...")

	if err := server.Shutdown(); err != nil {
		log.Fatal("Failed to shutdown server", zap.Error(err))
	}
}
