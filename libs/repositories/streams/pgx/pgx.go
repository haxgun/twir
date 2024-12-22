package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/repositories/streams/model"
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

var _ streams.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (model.Stream, error) {
	query := `
SELECT id,
       "userId",
       "userLogin",
       "userName",
       "gameId",
       "gameName",
       "communityIds",
       type,
       title,
       "viewerCount",
       "startedAt",
       language,
       "thumbnailUrl",
       "tagIds",
       "isMature",
       "parsedMessages",
       tags
FROM channels_streams
WHERE "userId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, streams.ErrNotFound
		}

		return model.Nil, fmt.Errorf("cannot get stream: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Stream])
	if err != nil {
		return model.Nil, fmt.Errorf("cannot scan stream: %w", err)
	}

	return result, nil
}
