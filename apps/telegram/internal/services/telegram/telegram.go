package telegram

import (
	"context"
	"errors"
	"strconv"

	channelsintegrationstelegramrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_telegram"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TelegramRepository channelsintegrationstelegramrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		telegramRepository: opts.TelegramRepository,
	}
}

type Service struct {
	telegramRepository channelsintegrationstelegramrepository.Repository
}

func (c *Service) CheckIsIntegrationExistsByChatID(ctx context.Context, chatID int64) (
	bool,
	error,
) {
	stringChatID := strconv.FormatInt(chatID, 10)

	_, err := c.telegramRepository.GetByChatID(ctx, stringChatID)
	if err != nil {
		if errors.Is(err, channelsintegrationstelegramrepository.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
