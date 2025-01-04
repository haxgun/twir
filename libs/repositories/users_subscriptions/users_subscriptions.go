package users_subscriptions

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/users_subscriptions/model"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (model.UserSubscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.UserSubscription, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateInput) (model.UserSubscription, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	UserID           string
	ExpireAt         time.Time
	SubscriptionID   uuid.UUID
	Provider         model.UserSubscriptionProvider
	TelegramChargeID *string
	ManualAssigned   bool
}

type UpdateInput struct {
	ExpireAt         *time.Time
	SubscriptionID   *uuid.UUID
	ManualAssigned   *bool
	TelegramChargeID *string
}
