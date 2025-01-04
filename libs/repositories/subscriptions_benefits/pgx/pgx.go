package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/subscriptions_benefits"
	"github.com/twirapp/twir/libs/repositories/subscriptions_benefits/model"
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

var _ subscriptions_benefits.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetMany(ctx context.Context) ([]model.SubscriptionBenefit, error) {
	query := `
SELECT id, name, created_at, updated_at, disable_inherit, quantity_value
FROM subscriptions_benefits
`
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	subscriptions, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.SubscriptionBenefit])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return subscriptions, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.SubscriptionBenefit, error) {
	query := `
SELECT id, name, created_at, updated_at, disable_inherit, quantity_value
FROM subscriptions_benefits
WHERE id = $1
LIMIT 1
	`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to execute query: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.SubscriptionBenefit],
	)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) GetByName(
	ctx context.Context,
	name model.SubscriptionBenefitName,
) (model.SubscriptionBenefit, error) {
	query := `
SELECT id, name, created_at, updated_at, disable_inherit, quantity_value
FROM subscriptions_benefits
WHERE name = $1
LIMIT 1
	`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, name)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to execute query: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.SubscriptionBenefit],
	)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input subscriptions_benefits.CreateInput,
) (model.SubscriptionBenefit, error) {
	query := `
INSERT INTO subscriptions_benefits (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at, disable_inherit, quantity_value
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.Name)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to execute query: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.SubscriptionBenefit],
	)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM subscriptions_benefits
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if rows.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete subscription benefit with id %s", id)
	}

	return nil
}
