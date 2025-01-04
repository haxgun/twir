-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE subscriptions_benefits (
	id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
	name       TEXT        NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX subscription_benefits_name_idx ON subscriptions_benefits (name);

CREATE TABLE subscription_tiers (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	price_cents INT NOT NULL,
	parent_id UUID references subscription_tiers(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX subscription_tiers_name_idx ON subscription_tiers (name);

CREATE TABLE subscription_tiers_benefits (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	subscription_id UUID NOT NULL REFERENCES subscription_tiers(id) ON DELETE CASCADE,
	benefit_id UUID NOT NULL REFERENCES subscriptions_benefits(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS subscription_tiers_benefits_subscription_id_idx ON subscription_tiers_benefits (subscription_id);

CREATE TYPE users_subscription_provider AS ENUM ('telegram');

CREATE TABLE IF NOT EXISTS users_subscriptions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expire_at TIMESTAMPTZ NOT NULL,
	subscription_id UUID NOT NULL REFERENCES subscription_tiers (id) ON DELETE CASCADE,
	manual_assigned BOOLEAN NOT NULL DEFAULT FALSE,
	provider users_subscription_provider NOT NULL,
	telegram_charge_id TEXT
);

CREATE INDEX IF NOT EXISTS users_subscriptions_user_id_idx ON users_subscriptions (user_id);

CREATE TABLE IF NOT EXISTS users_subscriptions_benefits (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	benefit_id UUID NOT NULL REFERENCES subscriptions_benefits(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expire_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS users_subscriptions_benefits_user_id_idx ON users_subscriptions_benefits (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
