package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTtsToGql, downTtsToGql)
}

type TtsToGqlEntity struct {
	ID        string
	Settings  []byte
	ChannelId string
	UserID    sql.NullString
}

type TtsToGqlCurrentSettings struct {
	Rate                               int      `json:"rate"`
	Pitch                              int      `json:"pitch"`
	Voice                              string   `json:"voice"`
	Volume                             int      `json:"volume"`
	Enabled                            bool     `json:"enabled"`
	MaxSymbols                         int      `json:"max_symbols"`
	DisallowedVoices                   []string `json:"disallowed_voices"`
	DoNotReadEmoji                     bool     `json:"do_not_read_emoji"`
	DoNotReadLinks                     bool     `json:"do_not_read_links"`
	ReadChatMessages                   bool     `json:"read_chat_messages"`
	DoNotReadTwitchEmotes              bool     `json:"do_not_read_twitch_emotes"`
	ReadChatMessagesNicknames          bool     `json:"read_chat_messages_nicknames"`
	AllowUsersChooseVoiceInMainCommand bool     `json:"allow_users_choose_voice_in_main_command"`
}

func upTtsToGql(ctx context.Context, tx *sql.Tx) error {
	voiceServiceQuery := `
CREATE TYPE channel_tts_voice_service AS ENUM (
	'rhvoice',
	'silero'
);
`
	if _, err := tx.Query(voiceServiceQuery); err != nil {
		return err
	}

	channelTableQuery := `
CREATE TABLE IF NOT EXISTS channels_overlays_tts (
	id uuid PRIMARY KEY default gen_random_uuid(),
	channel_id text references channels(id) ON DELETE CASCADE NOT NULL,
	enabled bool NOT NULL,
	rate int CHECK (rate >= 1 AND rate <= 100) NOT NULL,
	volume int CHECK (volume >= 1 AND volume <= 100) NOT NULL,
	pitch int CHECK (pitch >= 1 AND pitch <= 100) NOT NULL,
	voice varchar(255) NOT NULL,
	voice_service channel_tts_voice_service default 'rhvoice' NOT NULL,
	allow_users_choose_voice_in_main_command bool default false NOT NULL,
	max_symbols int NOT NULL,
	do_not_read_emoji bool DEFAULT true NOT NULL,
	do_not_read_twitch_emotes bool DEFAULT true NOT NULL,
	do_not_read_links bool DEFAULT true NOT NULL,
	read_chat_messages bool DEFAULT false NOT NULL,
	read_chat_messages_nicknames bool DEFAULT true NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_channels_overlays_tts_channel_id ON channels_overlays_tts (
channel_id);
`
	if _, err := tx.Query(channelTableQuery); err != nil {
		return err
	}

	disallowedVoicesTableQuery := `
CREATE TABLE IF NOT EXISTS channels_overlays_tts_disallowed_voices (
	id uuid PRIMARY KEY default gen_random_uuid(),
	channel_id text references channels_overlays_tts(channel_id) ON DELETE CASCADE NOT NULL,
	voice varchar(255) NOT NULL,
	service channel_tts_voice_service NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_channels_overlays_tts_disallowed_voices_channel_id ON channels_overlays_tts_disallowed_voices (channel_id);
`
	if _, err := tx.Query(disallowedVoicesTableQuery); err != nil {
		return err
	}

	userSettingsTableQuery := `
CREATE TABLE IF NOT EXISTS channels_overlays_tts_users_settings (
	id uuid PRIMARY KEY default gen_random_uuid(),
	channel_id text references channels_overlays_tts(channel_id) ON DELETE CASCADE NOT NULL,
	user_id text references users(id) ON DELETE CASCADE NOT NULL,
	voice varchar(255) NOT NULL,
	service channel_tts_voice_service NOT NULL,
	rate int CHECK (rate >= 1 AND rate <= 100) NOT NULL,
	pitch int CHECK (pitch >= 1 AND pitch <= 100) NOT NULL,
	volume int CHECK (volume >= 1 AND volume <= 100) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_channels_overlays_tts_users_settings_channel_id_user_id ON
channels_overlays_tts_users_settings (channel_id,user_id);
`
	if _, err := tx.Query(userSettingsTableQuery); err != nil {
		return err
	}

	var currentEntities []TtsToGqlEntity

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, settings, "channelId", "userId" FROM channels_modules_settings WHERE type = 'tts'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var entity TtsToGqlEntity
		if err := rows.Scan(
			&entity.ID,
			&entity.Settings,
			&entity.ChannelId,
			&entity.UserID,
		); err != nil {
			return err
		}
		currentEntities = append(currentEntities, entity)
	}

	for _, entity := range currentEntities {
		var settings TtsToGqlCurrentSettings
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return err
		}

		if settings.Rate == 0 {
			settings.Rate = 1
		}

		if settings.Pitch == 0 {
			settings.Pitch = 1
		}

		if !entity.UserID.Valid {
			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO channels_overlays_tts (
					id,
					channel_id,
					enabled,
					rate,
					volume,
					pitch,
					voice,
					allow_users_choose_voice_in_main_command,
					max_symbols,
					do_not_read_emoji,
					do_not_read_twitch_emotes,
					do_not_read_links,
					read_chat_messages,
					read_chat_messages_nicknames
				) VALUES (
					$1,
					$2,
					$3,
					$4,
					$5,
					$6,
					$7,
					$8,
					$9,
					$10,
					$11,
					$12,
					$13,
					$14
				)`,
				entity.ID,
				entity.ChannelId,
				settings.Enabled,
				settings.Rate,
				settings.Volume,
				settings.Pitch,
				settings.Voice,
				settings.AllowUsersChooseVoiceInMainCommand,
				settings.MaxSymbols,
				settings.DoNotReadEmoji,
				settings.DoNotReadTwitchEmotes,
				settings.DoNotReadLinks,
				settings.ReadChatMessages,
				settings.ReadChatMessagesNicknames,
			)
		} else {
			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO channels_overlays_tts_users_settings (
					id,
					channel_id,
					user_id,
					voice,
					service,
					rate,
					pitch,
					volume
				) VALUES (
					$1,
					$2,
					$3,
					$4,
					$5,
					$6,
					$7,
					$8
				)`,
				entity.ID,
				entity.ChannelId,
				entity.UserID.String,
				settings.Voice,
				"rhvoice",
				settings.Rate,
				settings.Pitch,
				settings.Volume,
			)

			if err != nil {
				return err
			}

			for _, disallowedVoice := range settings.DisallowedVoices {
				_, err = tx.ExecContext(
					ctx,
					`INSERT INTO channels_overlays_tts_disallowed_voices (
						id,
						channel_id,
						voice,
						service
					) VALUES (
						gen_random_uuid(),
						$1,
						$2,
						$3
					)`,
					entity.ChannelId,
					disallowedVoice,
					"rhvoice",
				)

				if err != nil {
					return err
				}
			}
		}

		_, err := tx.ExecContext(
			ctx,
			`DELETE FROM channels_modules_settings WHERE id = $1`,
			entity.ID,
		)
		if err != nil {
			return err
		}
	}

	// This code is executed when the migration is applied.
	return nil
}

func downTtsToGql(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
