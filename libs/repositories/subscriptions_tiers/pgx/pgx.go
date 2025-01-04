package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/subscriptions_tiers"
	"github.com/twirapp/twir/libs/repositories/subscriptions_tiers/model"
	"golang.org/x/exp/slices"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ subscriptions_tiers.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type scanModel struct {
	tier    model.SubscriptionTier
	benefit *model.SubscriptionTierBenefit
}

func (c *Pgx) scan(rows pgx.Rows) (scanModel, error) {
	tier := model.SubscriptionTier{}
	var benefitID, benefitSubscriptionID, benefitBenefitID sql.Null[uuid.UUID]
	var benefitCreatedAt, benefitUpdatedAt sql.Null[time.Time]

	result := scanModel{}

	err := rows.Scan(
		&tier.ID,
		&tier.Name,
		&tier.CreatedAt,
		&tier.UpdatedAt,
		&tier.PriceCents,
		&tier.ParentID,

		&benefitID,
		&benefitSubscriptionID,
		&benefitBenefitID,
		&benefitCreatedAt,
		&benefitUpdatedAt,
	)
	if err != nil {
		return scanModel{}, fmt.Errorf("failed to scan row: %w", err)
	}

	result.tier = tier

	if benefitID.Valid {
		result.benefit = &model.SubscriptionTierBenefit{
			ID:             benefitID.V,
			SubscriptionID: benefitSubscriptionID.V,
			BenefitID:      benefitBenefitID.V,
			CreatedAt:      benefitCreatedAt.V,
			UpdatedAt:      benefitUpdatedAt.V,
		}
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context) ([]model.SubscriptionTier, error) {
	query := `
SELECT st.id, st.name, st.created_at, st.updated_at, st.price_cents, st.parent_id,
       stb.id, stb.subscription_id, stb.benefit_id, stb.created_at, stb.updated_at
FROM subscription_tiers st
LEFT JOIN subscription_tiers_benefits stb ON st.id = stb.subscription_id
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	r, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer r.Close()

	tiersMap := make(map[uuid.UUID]*model.SubscriptionTier)
	for r.Next() {
		t, err := c.scan(r)
		if err != nil {
			return nil, err
		}

		if _, exists := tiersMap[t.tier.ID]; !exists {
			tiersMap[t.tier.ID] = &t.tier
		}

		if t.benefit != nil {
			tiersMap[t.tier.ID].Benefits = append(tiersMap[t.tier.ID].Benefits, *t.benefit)
		}
	}

	tiers := make([]model.SubscriptionTier, 0, len(tiersMap))
	for _, t := range tiersMap {
		tiers = append(tiers, *t)
	}

	slices.SortFunc(
		tiers, func(a, b model.SubscriptionTier) int {
			return strings.Compare(a.ID.String(), b.ID.String())
		},
	)

	return tiers, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.SubscriptionTier, error) {
	query := `
SELECT st.id, st.name, st.created_at, st.updated_at, st.price_cents, st.parent_id,
       stb.id, stb.subscription_id, stb.benefit_id, stb.created_at, stb.updated_at
FROM subscription_tiers st
LEFT JOIN subscription_tiers_benefits stb ON st.id = stb.subscription_id
WHERE st.id = $1
	`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tier model.SubscriptionTier
	var tierBenefits []model.SubscriptionTierBenefit

	for rows.Next() {
		t, err := c.scan(rows)
		if err != nil {
			return model.Nil, err
		}

		tier = t.tier
		if t.benefit != nil {
			tierBenefits = append(tierBenefits, *t.benefit)
		}
	}

	tier.Benefits = tierBenefits

	return tier, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input subscriptions_tiers.CreateInput,
) (model.SubscriptionTier, error) {
	query := `
INSERT INTO subscription_tiers (name, price_cents, parent_id)
VALUES ($1, $2, $3)
RETURNING id, name, created_at, updated_at, price_cents, parent_id
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.Name,
		input.PriceCents,
		input.ParentID,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to execute query: %w", err)
	}

	subscription, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.SubscriptionTier],
	)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect exactly one row: %w", err)
	}

	return subscription, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM subscription_tiers
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if rows.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete subscription tier: %w", pgx.ErrNoRows)
	}

	return nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input subscriptions_tiers.UpdateInput,
) (model.SubscriptionTier, error) {
	updateQuery := sq.Update("subscription_tiers").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, name, created_at, updated_at, price_cents, parent_id")
	updateQuery = repositories.SquirrelApplyPatch(
		updateQuery, map[string]interface{}{
			"name":        input.Name,
			"price_cents": input.PriceCents,
			"parent_id":   input.ParentID,
		},
	)

	query, args, err := updateQuery.ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to execute query: %w", err)
	}

	newTier, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.SubscriptionTier],
	)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect exactly one row: %w", err)
	}

	return c.GetByID(ctx, newTier.ID)
}

func (c *Pgx) AddBenefit(ctx context.Context, id uuid.UUID, benefitID uuid.UUID) error {
	query := `
INSERT INTO subscription_tiers_benefits (subscription_id, benefit_id)
VALUES ($1, $2)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id, benefitID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (c *Pgx) RemoveBenefit(ctx context.Context, id uuid.UUID, benefitID uuid.UUID) error {
	query := `
DELETE FROM subscription_tiers_benefits
WHERE subscription_id = $1 AND benefit_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id, benefitID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
