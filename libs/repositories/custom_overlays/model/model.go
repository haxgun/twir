package model

import (
	"time"

	"github.com/google/uuid"
)

type CustomOverlay struct {
	ID        uuid.UUID
	ChannelID string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Width     int
	Height    int
	Layers    []CustomOverlayLayer
}

type CustomOverlayLayer struct {
	ID                      uuid.UUID
	Type                    string
	OverlayID               uuid.UUID
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
