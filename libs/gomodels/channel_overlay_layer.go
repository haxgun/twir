package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type ChannelOverlayLayer struct {
	ID                      uuid.UUID          `gorm:"primary_key;column:id;type:UUID;"  json:"id"`
	Type                    ChannelOverlayType `gorm:"column:type;type:TEXT;"  json:"type"`
	OverlayID               uuid.UUID          `gorm:"column:overlay_id;type:UUID;"  json:"overlay_id"`
	Width                   int                `gorm:"column:width;type:INTEGER;"  json:"width"`
	Height                  int                `gorm:"column:height;type:INTEGER;"  json:"height"`
	CreatedAt               time.Time          `gorm:"column:created_at;data:timestamp;"  json:"createdAt"`
	UpdatedAt               time.Time          `gorm:"column:updated_at;data:timestamp;"  json:"updatedAt"`
	PeriodicallyRefetchData bool               `gorm:"column:periodically_refetch_data;type:BOOLEAN"  json:"periodically_refetch_data"`
	TransformString         string             `gorm:"column:transform_string;type:TEXT"  json:"transform_string"`

	// settings for html overlay
	SettingsHtmlHtml     *string `gorm:"column:settings_html_html;type:TEXT"  json:"settings_html_html"`
	SettingsHtmlCss      *string `gorm:"column:settings_html_css;type:TEXT"  json:"settings_html_css"`
	SettingsHtmlJs       *string `gorm:"column:settings_html_js;type:TEXT"  json:"settings_html_js"`
	SettingsHtmlDataPoll *int    `gorm:"column:settings_html_data_poll_seconds_interval;type:INTEGER"  json:"settings_html_data_poll"`

	// settings for image overlay
	SettingsImageURL *string `gorm:"column:settings_image_url;type:TEXT"  json:"settings_image_url"`

	Overlay *ChannelOverlay `gorm:"foreignKey:OverlayID" json:"overlay"`
}

func (c ChannelOverlayLayer) TableName() string {
	return "channels_overlays_layers"
}

// ChannelOverlayType types
type ChannelOverlayType string

func (e ChannelOverlayType) String() string {
	return string(e)
}

const (
	ChannelOverlayTypeHTML ChannelOverlayType = "HTML"
)

// ChannelOverlayLayerSettings settings
type ChannelOverlayLayerSettings struct {
	HtmlOverlayHTML                    string `json:"htmlOverlayHtml,omitempty"`
	HtmlOverlayCSS                     string `json:"htmlOverlayCss,omitempty"`
	HtmlOverlayJS                      string `json:"htmlOverlayJs,omitempty"`
	HtmlOverlayDataPollSecondsInterval int    `json:"htmlOverlayDataPollSecondsInterval,omitempty"`
}

func (a ChannelOverlayLayerSettings) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ChannelOverlayLayerSettings) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
