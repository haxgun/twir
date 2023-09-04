package processor

import (
	"context"

	"github.com/satont/twir/apps/dota-go/internal/cache"
	"github.com/satont/twir/apps/dota-go/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Processor) createOrUpdateAccount(
	ctx context.Context,
	channel cache.Channel,
	player *types.Player,
) error {
	if err := c.db.WithContext(ctx).Save(
		&model.ChannelsDotaAccounts{
			ID:        player.AccountID,
			ChannelID: channel.ID,
			NickName:  player.Name,
		},
	).Error; err != nil {
		return err
	}

	return nil
}
