package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Основной инстанс для базы. Неэкспортируемый.
// Будем использовать интерфейсы в каждом месте свои, чтобы ограничить поведение базы
// для конкретных мест
type postgres struct {
	pgx *pgxpool.Pool

	log *logger.ZapLogger

	config *pgxpool.Config
}

var _ service.UserRepository = (*postgres)(nil)

func New(cfg config.DatabaseConfig, logger *logger.ZapLogger) (*postgres, error) {
	pgxpoolConfig, err := pgxpoolConfig(cfg)
	if err != nil {
		logger.Error("[New()] cant parse config", zap.Error(err))
		return nil, err
	}

	return &postgres{
		config: pgxpoolConfig,
		log:    logger,
	}, nil
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

func (p *postgres) MustConnect(ctx context.Context) {
	// NewWithConfig creates a new Pool. config must have been created by [ParseConfig].
	// Поэтому нам его доставать через parseConfig (ф-я pgxpoolConfig)
	pool, err := pgxpool.NewWithConfig(ctx, p.config)
	if err != nil {
		p.log.Debug("[Connect()] error create pool with config", zap.Error(err))
		panic(err)
	}

	p.log.Debug("[Connect()] got pool", zap.Any("pgxpool", pool))

	if err := pool.Ping(ctx); err != nil {
		p.log.Debug("[Connect()] error Ping()", zap.Error(err))
		panic(err)
	}

	p.pgx = pool
}

// ------------------------------------------ SQL ------------------------------------------

func (p *postgres) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash
		FROM shop.users
		WHERE username = $1
	`

	row := p.pgx.QueryRow(ctx, query, username)

	user := &model.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *postgres) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	tx, err := p.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	getUserIDQuery := `
		SELECT id FROM shop.users WHERE username = $1
	`
	var toUserID uint
	err = tx.QueryRow(ctx, getUserIDQuery, params.ToUser).Scan(&toUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("recipient user not found")
		}
		return fmt.Errorf("failed to find recipient: %w", err)
	}

	lockBalanceQuery := `
		SELECT balance FROM shop.wallets WHERE user_id = $1 FOR UPDATE
	`
	var fromBalance int
	err = tx.QueryRow(ctx, lockBalanceQuery, params.FromUser).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("failed to fetch sender balance: %w", err)
	}

	if fromBalance < params.Amount {
		return fmt.Errorf("insufficient funds")
	}

	decreaseBalanceQuery := `
		UPDATE shop.wallets SET balance = balance - $1 WHERE user_id = $2
	`
	_, err = tx.Exec(ctx, decreaseBalanceQuery, params.Amount, params.FromUser)
	if err != nil {
		return fmt.Errorf("failed to debit sender: %w", err)
	}

	increaseBalanceQuery := `
		UPDATE shop.wallets SET balance = balance + $1 WHERE user_id = $2
	`
	_, err = tx.Exec(ctx, increaseBalanceQuery, params.Amount, toUserID)
	if err != nil {
		return fmt.Errorf("failed to credit recipient: %w", err)
	}

	insertTransactionQuery := `
		INSERT INTO shop.transactions (from_user_id, to_user_id, amount)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, insertTransactionQuery, params.FromUser, toUserID, params.Amount)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
