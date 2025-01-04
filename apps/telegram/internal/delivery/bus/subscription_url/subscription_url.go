package subscription_url

import (
	"context"

	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/telegram/internal/bot"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretelegram "github.com/twirapp/twir/libs/bus-core/telegram"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	TwirBus *buscore.Bus
	Logger  logger.Logger
	Bot     *bot.Bot
}

type Listener struct {
	bot *bot.Bot
}

func New(opts Opts) error {
	l := &Listener{
		bot: opts.Bot,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := opts.TwirBus.Telegram.GenerateSubscriptionUrl.SubscribeGroup(
					"telegram",
					l.GenerateSubscriptionUrl,
				); err != nil {
					return err
				}

				opts.Logger.Info("SubscriptionUrl listener started")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				opts.TwirBus.Telegram.GenerateSubscriptionUrl.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (c *Listener) GenerateSubscriptionUrl(
	ctx context.Context,
	data buscoretelegram.GenerateSubscriptionUrlInput,
) string {
	link := c.bot.CreateInvoiceLink(ctx, data.Title, data.StarsPrice)

	return link
}
