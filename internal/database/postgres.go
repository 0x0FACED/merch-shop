package database

import (
	"context"
	"fmt"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Postgres struct {
	pgx *pgxpool.Pool

	log *logger.ZapLogger

	config *pgxpool.Config
}

func New(cfg config.DatabaseConfig, logger *logger.ZapLogger) (*Postgres, error) {
	pgxpoolConfig, err := pgxpoolConfig(cfg)
	if err != nil {
		logger.Error("New() cant parse config", zap.Error(err))
		return nil, err
	}

	return &Postgres{
		config: pgxpoolConfig,
		log:    logger,
	}, nil
}

func (p Postgres) Pool() *pgxpool.Pool {
	return p.pgx
}

// pgxpoolConfig - функция, которая создает pgxpool.Config из конфига Database
// является неэкспортируемой, потому что используется только когда
// мы сздаем инстанс pgxpool
func pgxpoolConfig(cfg config.DatabaseConfig) (*pgxpool.Config, error) {
	// получаем структуру конфига из DSN
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// дополнительные настройки
	config.MaxConns = cfg.MaxOpenConns
	config.MinConns = cfg.MaxIdleConns
	config.MaxConnLifetime = cfg.ConnMaxLifetime
	config.MaxConnIdleTime = cfg.ConnMaxIdleTime

	config.ConnConfig.ConnectTimeout = cfg.ConnectionTimeout

	return config, nil
}

func (p *Postgres) MustConnect(ctx context.Context) {
	// NewWithConfig creates a new Pool. config must have been created by [ParseConfig].
	// Поэтому нам его доставать через parseConfig (ф-я pgxpoolConfig)
	pool, err := pgxpool.NewWithConfig(ctx, p.config)
	if err != nil {
		p.log.Debug("Connect() error create pool with config", zap.Error(err))
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		p.log.Debug("Connect() -> Ping() error", zap.Error(err))
		panic(err)
	}

	p.pgx = pool
}

func (p *Postgres) Close() error {
	p.log.Info("Closing database connection...")
	p.pgx.Close()
	p.log.Info("Database connection is closed")
	return nil
}
