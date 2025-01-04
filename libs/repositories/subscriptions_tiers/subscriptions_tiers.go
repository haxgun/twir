package subscriptions_tiers

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/subscriptions_tiers/model"
)

type Repository interface {
	GetMany(ctx context.Context) ([]model.SubscriptionTier, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.SubscriptionTier, error)
	Create(ctx context.Context, input CreateInput) (model.SubscriptionTier, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.SubscriptionTier, error)
	AddBenefit(ctx context.Context, id uuid.UUID, benefitID uuid.UUID) error
	RemoveBenefit(ctx context.Context, id uuid.UUID, benefitID uuid.UUID) error
}

type CreateInput struct {
	Name       string
	PriceCents int
	ParentID   *uuid.UUID
}

type UpdateInput struct {
	Name       *string
	PriceCents *int
	ParentID   *uuid.UUID
}
