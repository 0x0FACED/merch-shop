package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Основной инстанс для базы. Неэкспортируемый.
// Будем использовать интерфейсы в каждом месте свои, чтобы ограничить поведение базы
// для конкретных мест
type Postgres struct {
	pgx *pgxpool.Pool

	log *logger.ZapLogger

	config *pgxpool.Config
}

func New(cfg config.DatabaseConfig, logger *logger.ZapLogger) (*Postgres, error) {
	pgxpoolConfig, err := pgxpoolConfig(cfg)
	if err != nil {
		logger.Error("[New()] cant parse config", zap.Error(err))
		return nil, err
	}

	return &Postgres{
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

func (p *Postgres) MustConnect(ctx context.Context) {
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

// TODO:
// Надо еще создавать юзера при первой авторизации
// Надо делать либо так:
// 1. Отдаем ErrNotFound в сервис, сервис ее в API. Там проверяем.
// Если ErrNotFound вернулась, то вызываем метод для создания юзера
// 2. Сразу в базе создаем юзера
func (p *Postgres) AuthUser(ctx context.Context, params model.AuthUserParams) (*model.User, error) {
	query := `
		SELECT id, password_hash
		FROM shop.users
		WHERE username = $1
	`

	user := &model.User{}

	err := p.pgx.QueryRow(ctx, query, params.Username).Scan(&user.ID, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w query %q: %w", ErrQueryFailed, query, err)
	}

	return user, nil
}

func (p *Postgres) CreateUser(ctx context.Context, params model.CreateUserParams) (*model.User, error) {
	query := `
		INSERT INTO shop.users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username
	`

	user := &model.User{}

	err := p.pgx.QueryRow(ctx, query, params.Username, params.Password).Scan(
		&user.ID,
		&user.Username,
	)
	if err != nil {
		return nil, fmt.Errorf("%w query %q: %w", ErrQueryFailed, query, err)
	}

	createWalletQuery := `
		INSERT INTO shop.wallets (user_id)
		VALUES ($1)
	`

	// TODO: мб добавить rowsAffected првоерку у тега
	_, err = p.pgx.Exec(ctx, createWalletQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w query %q: %w", ErrQueryFailed, query, err)
	}

	return user, nil
}

func (p *Postgres) GetUserInfo(ctx context.Context, params model.GetUserInfoParams) (*model.UserInfo, error) {
	var balance uint
	err := p.pgx.QueryRow(ctx, `SELECT balance FROM shop.wallets WHERE user_id = $1`, params.ID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}

	inventoryQuery := `
		SELECT i.name, inv.quantity
		FROM shop.inventory inv
		JOIN shop.items i ON inv.item_id = i.id
		WHERE inv.user_id = $1
	`
	rows, err := p.pgx.Query(ctx, inventoryQuery, params.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}
	defer rows.Close()

	var inventory model.Inventory
	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrScanFailed, err)
		}
		items = append(items, item)
	}

	inventory.Items = items

	receivedQuery := `
		SELECT u.username, t.amount
		FROM shop.transactions t
		JOIN shop.users u ON t.from_user_id = u.id
		WHERE t.to_user_id = $1
	`
	rows, err = p.pgx.Query(ctx, receivedQuery, params.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}
	defer rows.Close()

	var received []model.ReceivedTransaction
	for rows.Next() {
		var trans model.ReceivedTransaction
		if err := rows.Scan(&trans.User, &trans.Amount); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrScanFailed, err)
		}
		received = append(received, trans)
	}

	sentQuery := `
		SELECT u.username, t.amount
		FROM shop.transactions t
		JOIN shop.users u ON t.to_user_id = u.id
		WHERE t.from_user_id = $1
	`
	rows, err = p.pgx.Query(ctx, sentQuery, params.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}
	defer rows.Close()

	var sent []model.SentTransaction
	for rows.Next() {
		var trans model.SentTransaction
		if err := rows.Scan(&trans.User, &trans.Amount); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrScanFailed, err)
		}
		sent = append(sent, trans)
	}

	return &model.UserInfo{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}

func (p *Postgres) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	p.log.Debug("SendCoin", zap.Any("params", params))

	tx, err := p.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToBeginTx, err)
	}

	defer func() {
		if err != nil {
			// TODO: мб error handle добавить для ролбэка
			tx.Rollback(ctx)
		}
	}()

	getUserIDQuery := `
		SELECT id FROM shop.users WHERE username = $1
	`
	var toUserID uint
	err = tx.QueryRow(ctx, getUserIDQuery, params.ToUser).Scan(&toUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("recipient %w", ErrNotFound)
		}
		return fmt.Errorf("%w query %q: %w", ErrFailedToFindRecipient, getUserIDQuery, err)
	}

	lockBalanceQuery := `
		SELECT balance FROM shop.wallets WHERE user_id = $1
	`
	var fromBalance int
	err = tx.QueryRow(ctx, lockBalanceQuery, params.FromUser).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("%w query %q: %w", ErrFailedToFetchBalance, lockBalanceQuery, err)
	}

	if fromBalance < params.Amount {
		return fmt.Errorf("%w", ErrInsufficientFunds)
	}

	decreaseBalanceQuery := `
		UPDATE shop.wallets SET balance = balance - $1 WHERE user_id = $2
	`
	_, err = tx.Exec(ctx, decreaseBalanceQuery, params.Amount, params.FromUser)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDebitSender, err)
	}

	increaseBalanceQuery := `
		UPDATE shop.wallets SET balance = balance + $1 WHERE user_id = $2
	`
	_, err = tx.Exec(ctx, increaseBalanceQuery, params.Amount, toUserID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreditRecipient, err)
	}

	insertTransactionQuery := `
		INSERT INTO shop.transactions (from_user_id, to_user_id, amount)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, insertTransactionQuery, params.FromUser, toUserID, params.Amount)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveTransaction, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCommitTx, err)
	}

	return nil
}
