package main

import (
	"github.com/twirapp/twir/apps/telegram/internal/bot"
	"github.com/twirapp/twir/apps/telegram/internal/bot/middlewares"
	"github.com/twirapp/twir/apps/telegram/internal/delivery/bus/subscription_url"
	"github.com/twirapp/twir/apps/telegram/internal/services/telegram"
	"github.com/twirapp/twir/libs/baseapp"
	"go.uber.org/fx"

	channelsintegrationstelegramrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_telegram"
	channelsintegrationstelegramrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_telegram/pgx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "telegram",
			},
		),
		fx.Provide(
			fx.Annotate(
				channelsintegrationstelegramrepositorypgx.NewFx,
				fx.As(new(channelsintegrationstelegramrepository.Repository)),
			),
		),
		fx.Provide(
			bot.New,
			telegram.New,
		),
		fx.Provide(
			middlewares.New,
		),
		fx.Invoke(
			// bot.New,
			subscription_url.New,
		),
		fx.NopLogger,
	).Run()
}
