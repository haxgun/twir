package model

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionBenefitName string

func (c SubscriptionBenefitName) String() string {
	return string(c)
}

const (
	SubscriptionBenefitNameCustomBot SubscriptionBenefitName = "CUSTOM_BOT"
)

type SubscriptionBenefit struct {
	ID             uuid.UUID
	Name           SubscriptionBenefitName
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DisableInherit bool
	QuantityValue  int
}

var Nil = SubscriptionBenefit{}
