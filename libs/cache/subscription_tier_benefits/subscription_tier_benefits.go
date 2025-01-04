package subscription_tier_benefits

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	benefitsrepository "github.com/twirapp/twir/libs/repositories/subscriptions_benefits"
	tiersrepository "github.com/twirapp/twir/libs/repositories/subscriptions_tiers"

	benefitsrepositorymodel "github.com/twirapp/twir/libs/repositories/subscriptions_benefits/model"
	tiersrepositorymodel "github.com/twirapp/twir/libs/repositories/subscriptions_tiers/model"
)

func New(
	tiersRepository tiersrepository.Repository,
	benefitsRepository benefitsrepository.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[[]benefitsrepositorymodel.SubscriptionBenefit] {
	return generic_cacher.New[[]benefitsrepositorymodel.SubscriptionBenefit](
		generic_cacher.Opts[[]benefitsrepositorymodel.SubscriptionBenefit]{
			Redis:     redis,
			KeyPrefix: "cache:twir:subscriptions_tiers_benefits:",
			LoadFn: func(ctx context.Context, key string) (
				[]benefitsrepositorymodel.SubscriptionBenefit,
				error,
			) {
				parsedTierKey, err := uuid.Parse(key)
				if err != nil {
					return nil, err
				}

				allTiers, err := tiersRepository.GetMany(ctx)
				if err != nil {
					return nil, err
				}

				allBenefits, err := benefitsRepository.GetMany(ctx)
				if err != nil {
					return nil, err
				}

				var foundTier tiersrepositorymodel.SubscriptionTier
				// find current tier by cache key
				for _, tier := range allTiers {
					if tier.ID != parsedTierKey {
						continue
					}

					foundTier = tier

					break
				}

				if foundTier.ID == uuid.Nil {
					return nil, fmt.Errorf("tier not found")
				}

				var currentTier *tiersrepositorymodel.SubscriptionTier
				currentTier = &foundTier

				var tierBenefits []benefitsrepositorymodel.SubscriptionBenefit
				for {
					if currentTier == nil {
						break
					}

					for _, benefit := range currentTier.Benefits {
						var foundBenefit *benefitsrepositorymodel.SubscriptionBenefit

						for _, b := range allBenefits {
							if b.ID == benefit.BenefitID {
								if b.DisableInherit && currentTier.ID != foundTier.ID {
									continue
								}

								foundBenefit = &b
								break
							}
						}

						if foundBenefit == nil {
							continue
						}

						tierBenefits = append(tierBenefits, *foundBenefit)
					}

					if currentTier.ParentID != nil && *currentTier.ParentID != uuid.Nil {
						for _, tier := range allTiers {
							if tier.ID == *currentTier.ParentID {
								currentTier = &tier
								break
							}
						}
					} else {
						currentTier = nil
					}
				}

				return tierBenefits, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
