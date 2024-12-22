package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channel_emote_usage"
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

var _ channel_emote_usage.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetUsedCount(ctx context.Context, input channel_emote_usage.GetUsedCount) (
	int,
	error,
) {
	query := sq.
		Select("COUNT(*)").
		From("channel_emote_usage").
		Where(squirrel.Eq{`"channel_id"`: input.ChannelID})

	if input.CreatedAtGte != nil {
		query = query.Where(`"channelId" >= ?`, input.CreatedAtGte)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build sql: %w", err)
	}

	var count int
	if err := c.pool.QueryRow(ctx, sql, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get used count: %w", err)
	}

	return count, nil
}
