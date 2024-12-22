package dashboard

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channel_emote_usage"
	"github.com/twirapp/twir/apps/api-gql/internal/services/streams"
	"github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	Logger                    logger.Logger
	StreamsService            *streams.Service
	TwitchService             *twitch.Service
	ChannelEmoteUsagesService *channel_emote_usage.Service
}

func New(opts Opts) *Service {
	return &Service{
		logger:                     opts.Logger,
		streamsService:             opts.StreamsService,
		twitchService:              opts.TwitchService,
		channelEmotesUsagesService: opts.ChannelEmoteUsagesService,
	}
}

type Service struct {
	logger                     logger.Logger
	streamsService             *streams.Service
	twitchService              *twitch.Service
	channelEmotesUsagesService *channel_emote_usage.Service
}

type GetDashboardStatsInput struct {
	ChannelID string

	WithFollowers bool
}

func (c *Service) GetDashboardStats(ctx context.Context, input GetDashboardStatsInput) (
	entity.DashboardStats,
	error,
) {
	if input.ChannelID == "" {
		return entity.DashboardStatsNil, fmt.Errorf("channel id is required")
	}

	stream, err := c.streamsService.GetByChannelID(ctx, input.ChannelID)
	if err != nil && !errors.Is(err, streamsrepository.ErrNotFound) {
		c.logger.Error("failed to get stream by channel id", slog.Any("err", err))
	}

	var wg errgroup.Group

	var parsedChatMessages int
	wg.Go(
		func() error {
			if stream.ID != "" {
				parsedChatMessages, err = c.streamsService.GetParsedChatMessages(ctx, stream.ID)
				if err != nil {
					c.logger.Error("failed to get parsed chat messages", slog.Any("err", err))
				}
			}

			return nil
		},
	)

	var followersCount int
	if input.WithFollowers {
		wg.Go(
			func() error {
				count, followersErr := c.twitchService.GetFollowersCountByChannelID(
					ctx,
					input.ChannelID,
				)
				if followersErr != nil {
					c.logger.Error("failed to get followers count", slog.Any("err", followersErr))
				}

				followersCount = count

				return nil
			},
		)
	}

	var channelInfo *helix.ChannelInformation
	wg.Go(
		func() error {
			info, infoErr := c.twitchService.GetChannelInformationByChannelID(ctx, input.ChannelID)
			if infoErr != nil {
				c.logger.Error("failed to get channel information", slog.Any("err", infoErr))
			}

			channelInfo = info
			return nil
		},
	)

	var usedEmotes int
	wg.Go(
		func() error {
			usageInput := channel_emote_usage.GetUsedCountInput{
				ChannelID:    input.ChannelID,
				CreatedAtGte: nil,
			}

			if stream.ID != "" {
				usageInput.CreatedAtGte = lo.ToPtr(stream.StartedAt)
			}

			count, err := c.channelEmotesUsagesService.GetUsedCount(ctx, usageInput)
			if err != nil {
				c.logger.Error("failed to get used emotes count", slog.Any("err", err))
			}

			usedEmotes = count

			return nil
		},
	)

	if err := wg.Wait(); err != nil {
		return entity.DashboardStatsNil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	stats := entity.DashboardStats{
		CategoryID:     "",
		CategoryName:   "",
		Viewers:        &stream.ViewerCount,
		StartedAt:      lo.ToPtr(stream.StartedAt),
		Title:          "",
		ChatMessages:   parsedChatMessages,
		Followers:      followersCount,
		UsedEmotes:     usedEmotes,
		RequestedSongs: 0,
		Subs:           0,
	}

	if channelInfo != nil {
		stats.CategoryName = channelInfo.GameName
		stats.Title = channelInfo.Title
		stats.CategoryID = channelInfo.GameID
	}

	return stats, nil
}
