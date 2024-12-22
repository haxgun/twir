package channel_emote_usage

import (
	"context"
	"time"
)

type Repository interface {
	GetUsedCount(ctx context.Context, input GetUsedCount) (int, error)
}

type GetUsedCount struct {
	ChannelID string

	CreatedAtGte *time.Time
}
