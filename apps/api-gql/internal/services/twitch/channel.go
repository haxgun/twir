package twitch

import (
	"context"

	"github.com/nicklaw5/helix/v2"
)

func (c *Service) GetChannelInformationByChannelID(
	ctx context.Context,
	channelID string,
) (*helix.ChannelInformation, error) {
	channelInformation, err := c.cachedTwitchClient.GetChannelInformationById(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return channelInformation, nil
}
