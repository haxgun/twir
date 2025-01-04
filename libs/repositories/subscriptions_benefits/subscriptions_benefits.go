package subscriptions_benefits

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/subscriptions_benefits/model"
)

type Repository interface {
	GetMany(ctx context.Context) ([]model.SubscriptionBenefit, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.SubscriptionBenefit, error)
	GetByName(ctx context.Context, name model.SubscriptionBenefitName) (
		model.SubscriptionBenefit,
		error,
	)
	Create(ctx context.Context, input CreateInput) (model.SubscriptionBenefit, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	Name model.SubscriptionBenefitName
}
