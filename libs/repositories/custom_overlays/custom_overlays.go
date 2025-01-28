package custom_overlays

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/custom_overlays/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.CustomOverlay, error)
	GetOne(ctx context.Context, id uuid.UUID) (model.CustomOverlay, error)
	Create(ctx context.Context, input CreateInput) (model.CustomOverlay, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.CustomOverlay, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type GetManyInput struct {
	ChannelID string
}

type LayerInput struct {
	Type                    string
	Width                   int
	Height                  int
	CreatedAt               time.Time
	UpdatedAt               time.Time
	PeriodicallyRefetchData bool
	TransformString         string

	// settings for html overlay
	SettingsHtmlHtml     *string
	SettingsHtmlCss      *string
	SettingsHtmlJs       *string
	SettingsHtmlDataPoll *int

	// settings for image overlay
	SettingsImageURL *string
}

type CreateInput struct {
	ChannelID string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Width     int
	Height    int
	Layers    []LayerInput
}

type UpdateInput struct {
	ChannelID string
	Name      string
	Width     int
	Height    int
	Layers    []LayerInput
}
