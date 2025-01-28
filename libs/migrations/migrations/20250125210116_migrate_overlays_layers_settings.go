package migrations

import (
	"context"
	"database/sql"

	"github.com/goccy/go-json"
	"github.com/pressly/goose/v3"
)

type channelOverlayLayer20250125210116 struct {
	ID       string `json:"id"`
	Settings []byte `json:"settings"`
}

type channelOverlayLayerSettings20250125210116 struct {
	HtmlOverlayHTML                    string `json:"htmlOverlayHtml,omitempty"`
	HtmlOverlayCSS                     string `json:"htmlOverlayCss,omitempty"`
	HtmlOverlayJS                      string `json:"htmlOverlayJs,omitempty"`
	HtmlOverlayDataPollSecondsInterval int    `json:"htmlOverlayDataPollSecondsInterval,omitempty"`
}

func init() {
	goose.AddMigrationContext(upMigrateOverlaysLayersSettings, downMigrateOverlaysLayersSettings)
}

func upMigrateOverlaysLayersSettings(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, settings FROM channels_overlays_layers`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var entities []channelOverlayLayer20250125210116

	for rows.Next() {
		var entity channelOverlayLayer20250125210116
		if err := rows.Scan(&entity.ID, &entity.Settings); err != nil {
			return err
		}
		entities = append(entities, entity)
	}

	_, err = tx.ExecContext(
		ctx,
		`ALTER TABLE channels_overlays_layers DROP COLUMN settings`,
	)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
ALTER TABLE channels_overlays_layers ADD COLUMN settings_html_html TEXT;
ALTER TABLE channels_overlays_layers ADD COLUMN settings_html_css	TEXT;
ALTER TABLE channels_overlays_layers ADD COLUMN settings_html_js TEXT;
ALTER TABLE channels_overlays_layers ADD COLUMN settings_html_data_poll_seconds_interval INT;
`,
	)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		var settings channelOverlayLayerSettings20250125210116
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`UPDATE channels_overlays_layers SET settings_html_html = $1, settings_html_css = $2, settings_html_js = $3, settings_html_data_poll_seconds_interval = $4 WHERE id = $5`,
			settings.HtmlOverlayHTML,
			settings.HtmlOverlayCSS,
			settings.HtmlOverlayJS,
			settings.HtmlOverlayDataPollSecondsInterval,
			entity.ID,
		)
	}

	return nil
}

func downMigrateOverlaysLayersSettings(ctx context.Context, tx *sql.Tx) error {
	return nil
}
