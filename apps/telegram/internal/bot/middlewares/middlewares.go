package middlewares

import (
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/telegram/internal/services/telegram"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Logger          logger.Logger
	TelegramService *telegram.Service
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		logger:          opts.Logger,
		telegramService: opts.TelegramService,
	}
}

type Middlewares struct {
	logger          logger.Logger
	telegramService *telegram.Service
}
