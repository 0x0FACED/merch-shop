package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testServer *server.Server

func TestMain(m *testing.M) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoadTestConfig()
	var err error
	testServer, err = server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("testServer created:", testServer)

	clearDB(ctx, testServer.Database())

	go func() {
		if err := testServer.Start(ctx); err != nil {
			log.Fatalf("Server stopped with err: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	code := m.Run()

	if err := testServer.Shutdown(); err != nil {
		log.Fatal("Failed to shutdown test server:", err)
	}

	os.Exit(code)
}

func clearDB(ctx context.Context, db *pgxpool.Pool) {
	db.Exec(ctx, "DELETE FROM shop.users")
	db.Exec(ctx, "DELETE FROM shop.wallets")
	db.Exec(ctx, "DELETE FROM shop.inventory")
	db.Exec(ctx, "DELETE FROM shop.transactions")
}
