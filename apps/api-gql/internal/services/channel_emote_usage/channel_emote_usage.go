package channel_emote_usage

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/channel_emote_usage"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelEmoteUsageRepository channel_emote_usage.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repository: opts.ChannelEmoteUsageRepository,
	}
}

type Service struct {
	repository channel_emote_usage.Repository
}

type GetUsedCountInput struct {
	ChannelID string

	CreatedAtGte *time.Time
}

func (c *Service) GetUsedCount(ctx context.Context, input GetUsedCountInput) (
	int,
	error,
) {
	count, err := c.repository.GetUsedCount(
		ctx, channel_emote_usage.GetUsedCount{
			ChannelID:    input.ChannelID,
			CreatedAtGte: input.CreatedAtGte,
		},
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}
