package integrations

import (
	config "github.com/satont/twir/libs/config"
	channelsintegrationstelegramrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_telegram"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config             config.Config
	TelegramRepository channelsintegrationstelegramrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		config:             opts.Config,
		telegramRepository: opts.TelegramRepository,
	}
}

type Service struct {
	config             config.Config
	telegramRepository channelsintegrationstelegramrepository.Repository
}
