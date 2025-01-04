package model

import (
	"github.com/google/uuid"
)

type TelegramIntegration struct {
	ID             uuid.UUID
	ChannelID      string
	TelegramChatID string
	CreatedAt      string
	UpdatedAt      string
}
