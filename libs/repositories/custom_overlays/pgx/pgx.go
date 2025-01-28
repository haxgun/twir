package pgx

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/custom_overlays"
	"github.com/twirapp/twir/libs/repositories/custom_overlays/model"
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

var _ custom_overlays.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"overlay.id",
	"overlay.channel_id",
	"overlay.name",
	"overlay.created_at",
	"overlay.updated_at",
	"overlay.width",
	"overlay.height",
	"layer.id",
	"layer.type",
	"layer.overlay_id",
	"layer.width",
	"layer.height",
	"layer.created_at",
	"layer.updated_at",
	"layer.periodically_refetch_data",
	"layer.transform_string",
	"layer.settings_html_html",
	"layer.settings_html_css",
	"layer.settings_html_js",
	"layer.settings_html_data_poll_seconds_interval",
	"layer.settings_image_url",
}

type scanModel struct {
	overlay model.CustomOverlay
	layer   model.CustomOverlayLayer
}

func (c *Pgx) scan(rows pgx.Row) (scanModel, error) {
	var m scanModel
	err := rows.Scan(
		&m.overlay.ID,
		&m.overlay.ChannelID,
		&m.overlay.Name,
		&m.overlay.CreatedAt,
		&m.overlay.UpdatedAt,
		&m.overlay.Width,
		&m.overlay.Height,
		&m.layer.ID,
		&m.layer.Type,
		&m.layer.OverlayID,
		&m.layer.Width,
		&m.layer.Height,
		&m.layer.CreatedAt,
		&m.layer.UpdatedAt,
		&m.layer.PeriodicallyRefetchData,
		&m.layer.TransformString,
		&m.layer.SettingsHtmlHtml,
		&m.layer.SettingsHtmlCss,
		&m.layer.SettingsHtmlJs,
		&m.layer.SettingsHtmlDataPoll,
		&m.layer.SettingsImageURL,
	)

	if err != nil {
		return scanModel{}, fmt.Errorf("error while scanning row: %w", err)
	}

	return m, nil
}

func (c *Pgx) GetOne(ctx context.Context, id uuid.UUID) (model.CustomOverlay, error) {
	query, args, err := sq.
		Select(selectColumns...).
		From("channels_overlays overlay").
		Join("channels_overlays_layers layer ON overlay.id = layer.overlay_id").
		Where(squirrel.Eq{"overlay.id": id}).
		ToSql()
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while building query: %w", err)
	}

	overlay := model.CustomOverlay{}
	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while querying overlay: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		m, scanErr := c.scan(rows)
		if scanErr != nil {
			return model.CustomOverlay{}, fmt.Errorf("error while scanning row: %w", scanErr)
		}

		if overlay.ID == uuid.Nil {
			overlay = m.overlay
		}

		overlay.Layers = append(overlay.Layers, m.layer)
	}

	if rows.Err() != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while scanning rows: %w", rows.Err())
	}

	return overlay, nil
}

func (c *Pgx) GetMany(
	ctx context.Context,
	input custom_overlays.GetManyInput,
) ([]model.CustomOverlay, error) {
	query, args, err := sq.
		Select(selectColumns...).
		From("channels_overlays overlay").
		Join("channels_overlays_layers layer ON overlay.id = layer.overlay_id").
		Where(squirrel.Eq{"overlay.channel_id": input.ChannelID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error while building query: %w", err)
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error while querying overlays: %w", err)
	}

	defer rows.Close()

	overlays := map[uuid.UUID]*model.CustomOverlay{}

	for rows.Next() {
		m, scanErr := c.scan(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("error while scanning row: %w", scanErr)
		}

		if overlay, ok := overlays[m.overlay.ID]; ok {
			overlay.Layers = append(overlay.Layers, m.layer)
		} else {
			overlays[m.overlay.ID] = &m.overlay
			overlays[m.overlay.ID].Layers = append(overlays[m.overlay.ID].Layers, m.layer)
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while scanning rows: %w", rows.Err())
	}

	overlaysSlice := make([]model.CustomOverlay, 0, len(overlays))
	for _, overlay := range overlays {
		overlaysSlice = append(overlaysSlice, *overlay)
	}

	slices.SortFunc(
		overlaysSlice, func(a, b model.CustomOverlay) int {
			return cmp.Compare(a.CreatedAt.UnixMilli(), b.CreatedAt.UnixMilli())
		},
	)

	return overlaysSlice, nil
}

func (c *Pgx) Create(ctx context.Context, input custom_overlays.CreateInput) (
	model.CustomOverlay,
	error,
) {
	overlayQuery, overlayArgs, err := sq.
		Insert("channels_overlays").
		Columns("channel_id", "name", "width", "height").
		Values(input.ChannelID, input.Name, input.Width, input.Height).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while building query: %w", err)
	}

	var overlayID uuid.UUID
	tx, err := c.pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while starting transaction: %w", err)
	}

	err = tx.QueryRow(ctx, overlayQuery, overlayArgs...).Scan(&overlayID)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while inserting overlay: %w", err)
	}

	for _, layer := range input.Layers {
		layerQuery, layerArgs, err := sq.
			Insert("channels_overlays_layers").
			Columns(
				"type",
				"overlay_id",
				"width",
				"height",
				"periodically_refetch_data",
				"transform_string",
				"settings_html_html",
				"settings_html_css",
				"settings_html_js",
				"settings_html_data_poll_seconds_interval",
				"settings_image_url",
			).
			Values(
				layer.Type,
				overlayID,
				layer.Width,
				layer.Height,
				layer.PeriodicallyRefetchData,
				layer.TransformString,
				layer.SettingsHtmlHtml,
				layer.SettingsHtmlCss,
				layer.SettingsHtmlJs,
				layer.SettingsHtmlDataPoll,
				layer.SettingsImageURL,
			).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return model.CustomOverlay{}, fmt.Errorf("error while building layer query: %w", err)
		}

		var layerID uuid.UUID
		err = tx.QueryRow(ctx, layerQuery, layerArgs...).Scan(&layerID)
		if err != nil {
			return model.CustomOverlay{}, fmt.Errorf("error while inserting layer: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while committing transaction: %w", err)
	}

	return c.GetOne(ctx, overlayID)
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input custom_overlays.UpdateInput) (
	model.CustomOverlay,
	error,
) {
	overlayQuery, overlayArgs, err := sq.
		Update("channels_overlays").
		Set("name", input.Name).
		Set("width", input.Width).
		Set("height", input.Height).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while building query: %w", err)
	}

	var overlayID uuid.UUID
	tx, err := c.pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while starting transaction: %w", err)
	}

	err = tx.QueryRow(ctx, overlayQuery, overlayArgs...).Scan(&overlayID)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while updating overlay: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM channels_overlays_layers WHERE overlay_id = $1", overlayID)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while deleting layers: %w", err)
	}

	for _, layer := range input.Layers {
		layerQuery, layerArgs, err := sq.
			Insert("channels_overlays_layers").
			Columns(
				"type",
				"overlay_id",
				"width",
				"height",
				"periodically_refetch_data",
				"transform_string",
				"settings_html_html",
				"settings_html_css",
				"settings_html_js",
				"settings_html_data_poll_seconds_interval",
				"settings_image_url",
			).
			Values(
				layer.Type,
				overlayID,
				layer.Width,
				layer.Height,
				layer.PeriodicallyRefetchData,
				layer.TransformString,
				layer.SettingsHtmlHtml,
				layer.SettingsHtmlCss,
				layer.SettingsHtmlJs,
				layer.SettingsHtmlDataPoll,
				layer.SettingsImageURL,
			).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return model.CustomOverlay{}, fmt.Errorf("error while building layer query: %w", err)
		}

		var layerID uuid.UUID
		err = tx.QueryRow(ctx, layerQuery, layerArgs...).Scan(&layerID)
		if err != nil {
			return model.CustomOverlay{}, fmt.Errorf("error while updating layer: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.CustomOverlay{}, fmt.Errorf("error while committing transaction: %w", err)
	}

	return c.GetOne(ctx, overlayID)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM channels_overlays WHERE id = $1"
	_, err := c.pool.Exec(ctx, query, id)
	return fmt.Errorf("error while deleting overlay: %w", err)
}
