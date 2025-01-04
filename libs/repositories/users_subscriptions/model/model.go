package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSubscriptionProvider string

const (
	UserSubscriptionProviderTelegram UserSubscriptionProvider = "telegram"
)

type UserSubscription struct {
	ID               uuid.UUID
	UserID           string
	CreatedAt        time.Time
	ExpireAt         time.Time
	SubscriptionID   uuid.UUID
	ManualAssigned   bool
	Provider         UserSubscriptionProvider
	TelegramChargeID *string
}

var Nil = UserSubscription{}
