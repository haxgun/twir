package channels_integrations_telegram

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_integrations_telegram/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.TelegramIntegration, error)
	GetByChatID(ctx context.Context, chatID string) (model.TelegramIntegration, error)
	Create(ctx context.Context, input CreateInput) (model.TelegramIntegration, error)
}

type CreateInput struct {
	ChannelID      string
	TelegramChatID string
}
