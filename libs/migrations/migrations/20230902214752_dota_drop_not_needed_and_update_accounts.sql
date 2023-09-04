-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE "dota_matches_cards";
DROP TABLE "dota_matches_results";
DROP TABLE "dota_matches";
DROP TABLE "dota_game_modes";
DROP TABLE "dota_heroes";
DELETE FROM "channels_dota_accounts";
ALTER TABLE "channels_dota_accounts" ADD COLUMN "nickname" text NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
