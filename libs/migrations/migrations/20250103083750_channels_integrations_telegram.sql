-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_integrations_telegram" (
		"id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		"channel_id" TEXT NOT NULL REFERENCES "channels" ("id") ON DELETE CASCADE,
		"telegram_chat_id" TEXT NOT NULL,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX "channels_integrations_telegram_channel_id" ON "channels_integrations_telegram" ("channel_id");
CREATE UNIQUE INDEX "channels_integrations_telegram_telegram_chat_id" ON "channels_integrations_telegram" ("telegram_chat_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
