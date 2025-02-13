package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/jackc/pgx/v5"
)

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

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrScanFailed, err)
		}
		items = append(items, item)
	}

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
		Inventory: items,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}

func (p *Postgres) SendCoin(ctx context.Context, params model.SendCoinParams) error {
	tx, err := p.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToBeginTx, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
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
		return ErrInsufficientFunds
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

func (p *Postgres) GetUserBalance(ctx context.Context, userID uint) (uint, error) {
	var balance uint
	err := p.pgx.QueryRow(ctx, `SELECT balance FROM shop.wallets WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}

	return balance, nil
}

func (p *Postgres) BuyItem(ctx context.Context, params model.BuyItemParams) error {
	tx, err := p.pgx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToBeginTx, err)
	}
	defer tx.Rollback(ctx)

	var itemID, price uint
	err = tx.QueryRow(ctx, `SELECT id, price FROM shop.items WHERE name = $1`, params.Item).Scan(&itemID, &price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}

	if params.Balance < price {
		return ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, `UPDATE shop.wallets SET balance = balance - $1 WHERE user_id = $2`, price, params.UserID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO shop.inventory (user_id, item_id, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, item_id) DO UPDATE
		SET quantity = shop.inventory.quantity + 1
	`, params.UserID, itemID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrQueryFailed, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCommitTx, err)
	}

	return nil
}
