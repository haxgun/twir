package model

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionTier struct {
	ID         uuid.UUID
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	PriceCents int
	ParentID   *uuid.UUID

	Benefits []SubscriptionTierBenefit
}

var Nil = SubscriptionTier{}

type SubscriptionTierBenefit struct {
	ID             uuid.UUID
	SubscriptionID uuid.UUID
	BenefitID      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
