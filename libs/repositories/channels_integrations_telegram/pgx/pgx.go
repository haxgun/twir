package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_telegram"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_telegram/model"
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

var _ channels_integrations_telegram.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChatID(ctx context.Context, chatID string) (model.TelegramIntegration, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	query := `
SELECT "id", "channel_id", "telegram_chat_id", "created_at", "updated_at"
FROM "channels_integrations_telegram"
WHERE "telegram_chat_id" = $1
LIMIT 1
`

	rows, err := conn.Query(ctx, query, chatID)
	if err != nil {
		return model.TelegramIntegration{}, fmt.Errorf(
			"failed to get telegram integration by chat id: %w",
			err,
		)
	}

	integration, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.TelegramIntegration],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.TelegramIntegration{}, channels_integrations_telegram.ErrNotFound
		}

		return model.TelegramIntegration{}, fmt.Errorf(
			"failed to parse telegram integration by chat id: %w",
			err,
		)
	}

	return integration, nil
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.TelegramIntegration,
	error,
) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	query := `
SELECT "id", "channel_id", "telegram_chat_id", "created_at", "updated_at"
FROM "channels_integrations_telegram"
WHERE "channel_id" = $1
LIMIT 1
`

	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.TelegramIntegration{}, fmt.Errorf(
			"failed to get telegram integration by channel id: %w",
			err,
		)
	}

	integration, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.TelegramIntegration],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.TelegramIntegration{}, channels_integrations_telegram.ErrNotFound
		}

		return model.TelegramIntegration{}, fmt.Errorf(
			"failed to parse telegram integration by channel id: %w",
			err,
		)
	}

	return integration, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input channels_integrations_telegram.CreateInput,
) (model.TelegramIntegration, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	query := `
INSERT INTO "channels_integrations_telegram" ("channel_id", "telegram_chat_id")
VALUES ($1, $2)
`

	_, err := conn.Exec(
		ctx,
		query,
		input.ChannelID,
		input.TelegramChatID,
	)
	if err != nil {
		return model.TelegramIntegration{}, fmt.Errorf("failed to create telegram integration: %w", err)
	}

	return c.GetByChannelID(ctx, input.ChannelID)
}
