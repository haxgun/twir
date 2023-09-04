package processor

import (
	"context"

	"github.com/satont/twir/apps/dota-go/internal/cache"
	"github.com/satont/twir/apps/dota-go/internal/types"
)

func (c *Processor) processMatch(
	ctx context.Context,
	channel cache.Channel,
	player *types.Player,
	mapData *types.MapData,
) error {
	isExists, err := c.cache.MatchIsExists(ctx, channel.ID, player.AccountID, mapData.MatchID)
	if err != nil {
		return err
	}

	if isExists {

	}

	return nil
}
