package processor

import (
	"context"

	"github.com/paralin/go-dota2"
	"github.com/satont/twir/apps/dota-go/internal/cache"
	"github.com/satont/twir/apps/dota-go/internal/types"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Processor struct {
	db     *gorm.DB
	cache  *cache.Cache
	logger *logrus.Logger
	dota   *dota2.Dota2
}

type Opts struct {
	fx.In

	Db     *gorm.DB
	Cache  *cache.Cache
	Logger *logrus.Logger
	Dota   *dota2.Dota2
}

func New(opts Opts) *Processor {
	return &Processor{
		db:     opts.Db,
		cache:  opts.Cache,
		logger: opts.Logger,
		dota:   opts.Dota,
	}
}

func (c *Processor) Process(ctx context.Context, packet *types.Packet) error {
	channel, err := c.cache.ChannelByToken(ctx, packet.Auth.Token)
	if err != nil {
		return err
	}

	if packet.Player != nil && packet.Player.AccountID != "" {
		if err := c.createOrUpdateAccount(ctx, channel, packet.Player); err != nil {
			return err
		}

		if packet.Map != nil && packet.Map.MatchID != "" {
			if err := c.processMatch(ctx, channel, packet.Player, packet.Map); err != nil {
				return err
			}
		}
	}

	return nil
}
