package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/users_subscriptions"
	"github.com/twirapp/twir/libs/repositories/users_subscriptions/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ users_subscriptions.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByUserID(ctx context.Context, userID string) (model.UserSubscription, error) {
	query := `
SELECT id, user_id, created_at, expire_at, subscription_id, manual_assigned, provider, telegram_charge_id
FROM users_subscriptions
WHERE user_id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		return model.Nil, fmt.Errorf("cannot query: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserSubscription])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, users_subscriptions.ErrNotFound
		}

		return model.Nil, fmt.Errorf("cannot collect exactly one row: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.UserSubscription, error) {
	query := `
SELECT id, user_id, created_at, expire_at, subscription_id, manual_assigned, provider, telegram_charge_id
FROM users_subscriptions
WHERE id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.Nil, fmt.Errorf("cannot query: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserSubscription])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, users_subscriptions.ErrNotFound
		}

		return model.Nil, fmt.Errorf("cannot collect exactly one row: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM users_subscriptions"

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	err := conn.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("cannot query row: %w", err)
	}

	return count, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input users_subscriptions.CreateInput,
) (model.UserSubscription, error) {
	query := `
INSERT INTO users_subscriptions (user_id, expire_at, subscription_id, manual_assigned, provider, telegram_charge_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, created_at, expire_at, subscription_id, manual_assigned, provider, telegram_charge_id
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.UserID,
		input.ExpireAt,
		input.SubscriptionID,
		input.ManualAssigned,
		input.Provider,
		input.TelegramChargeID,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("cannot exec: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserSubscription])
	if err != nil {
		return model.Nil, fmt.Errorf("cannot collect exactly one row: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input users_subscriptions.UpdateInput,
) error {
	updateBuilder := sq.Update("users_subscriptions").
		Where(squirrel.Eq{"id": id})
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder, map[string]interface{}{
			"expire_at":          input.ExpireAt,
			"subscription_id":    input.SubscriptionID,
			"manual_assigned":    input.ManualAssigned,
			"telegram_charge_id": input.TelegramChargeID,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("cannot build sql: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot exec: %w", err)
	}

	return nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users_subscriptions WHERE id = $1"

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot exec: %w", err)
	}

	return nil
}
