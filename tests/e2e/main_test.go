package e2e

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/server"
	"github.com/0x0FACED/merch-shop/internal/server/handler"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	testServer *server.Server
	testDB     *database.Postgres
)

func TestMain(m *testing.M) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoadTestConfig()
	log := logger.New(cfg.Logger)

	var err error
	testDB, err = database.New(cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to create test database", zap.Error(err))
	}
	testDB.MustConnect(ctx)
	defer testDB.Close()

	merchService := service.NewUserService(testDB, log)
	h := handler.NewHandler(merchService, log, &cfg.Server)

	testServer, err = server.NewServer(cfg, h)
	if err != nil {
		log.Fatal("Failed to start test server", zap.Error(err))
	}

	clearDB(ctx, testDB.Pool())

	go func() {
		if err := testServer.Start(ctx); err != nil {
			log.Fatal("Server stopped with err", zap.Error(err))
		}
	}()

	time.Sleep(2 * time.Second)

	code := m.Run()

	if err := testServer.Shutdown(); err != nil {
		log.Fatal("Failed to shutdown test server", zap.Error(err))
	}

	os.Exit(code)
}

func clearDB(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DELETE FROM shop.users")
	_, _ = db.Exec(ctx, "DELETE FROM shop.wallets")
	_, _ = db.Exec(ctx, "DELETE FROM shop.inventory")
	_, _ = db.Exec(ctx, "DELETE FROM shop.transactions")
}
