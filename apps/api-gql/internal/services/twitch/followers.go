package twitch

import (
	"context"
)

func (c *Service) GetFollowersCountByChannelID(ctx context.Context, channelID string) (int, error) {
	return c.cachedTwitchClient.GetChannelFollowersCountByChannelId(
		ctx,
		channelID,
	)
}
