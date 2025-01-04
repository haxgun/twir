-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE subscriptions_benefits ADD COLUMN disable_inherit BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE subscriptions_benefits ADD COLUMN quantity_value INT NOT NULL DEFAULT 0;
DROP INDEX subscription_benefits_name_idx;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
