package custom_overlays

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"

	customoverlaysrepository "github.com/twirapp/twir/libs/repositories/custom_overlays"
	customoverlaysrepositorymodel "github.com/twirapp/twir/libs/repositories/custom_overlays/model"
)

type Opts struct {
	fx.In

	Repository     customoverlaysrepository.Repository
	TwirBus        *buscore.Bus
	WebsocketsGrpc websockets.WebsocketClient
}

func New(opts Opts) *Service {
	return &Service{
		repository:     opts.Repository,
		twirBus:        opts.TwirBus,
		websocketsGrpc: opts.WebsocketsGrpc,
	}
}

type Service struct {
	repository     customoverlaysrepository.Repository
	twirBus        *buscore.Bus
	websocketsGrpc websockets.WebsocketClient
}

func (c *Service) mapToEntity(m customoverlaysrepositorymodel.CustomOverlay) entity.CustomOverlay {
	layers := make([]entity.CustomOverlayLayer, len(m.Layers))
	for i, layer := range m.Layers {
		layers[i] = entity.CustomOverlayLayer{
			ID:                      layer.ID,
			Type:                    entity.ChannelOverlayLayerType(layer.Type),
			OverlayID:               layer.OverlayID,
			Width:                   layer.Width,
			Height:                  layer.Height,
			CreatedAt:               layer.CreatedAt,
			UpdatedAt:               layer.UpdatedAt,
			PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
			TransformString:         layer.TransformString,
			SettingsHtmlHtml:        layer.SettingsHtmlHtml,
			SettingsHtmlCss:         layer.SettingsHtmlCss,
			SettingsHtmlJs:          layer.SettingsHtmlJs,
			SettingsHtmlDataPoll:    layer.SettingsHtmlDataPoll,
			SettingsImageURL:        layer.SettingsImageURL,
		}
	}

	return entity.CustomOverlay{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Width:     m.Width,
		Height:    m.Height,
		Layers:    layers,
	}
}

func (c *Service) GetOne(ctx context.Context, id uuid.UUID) (entity.CustomOverlay, error) {
	overlay, err := c.repository.GetOne(ctx, id)
	if err != nil {
		return entity.CustomOverlay{}, err
	}

	return c.mapToEntity(overlay), nil
}

type GetManyInput struct {
	ChannelID string
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.CustomOverlay, error) {
	overlays, err := c.repository.GetMany(
		ctx, customoverlaysrepository.GetManyInput{
			ChannelID: input.ChannelID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.CustomOverlay, len(overlays))
	for i, overlay := range overlays {
		result[i] = c.mapToEntity(overlay)
	}

	return result, nil
}

type CreateInput struct {
	ChannelID string
	Name      string
	Width     int
	Height    int
	Layers    []LayerInput
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

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.CustomOverlay, error) {
	layers := make([]customoverlaysrepository.LayerInput, len(input.Layers))
	for i, layer := range input.Layers {
		layers[i] = customoverlaysrepository.LayerInput{
			Type:                    layer.Type,
			Width:                   layer.Width,
			Height:                  layer.Height,
			CreatedAt:               layer.CreatedAt,
			UpdatedAt:               layer.UpdatedAt,
			PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
			TransformString:         layer.TransformString,
			SettingsHtmlHtml:        layer.SettingsHtmlHtml,
			SettingsHtmlCss:         layer.SettingsHtmlCss,
			SettingsHtmlJs:          layer.SettingsHtmlJs,
			SettingsHtmlDataPoll:    layer.SettingsHtmlDataPoll,
			SettingsImageURL:        layer.SettingsImageURL,
		}
	}

	overlay, err := c.repository.Create(
		ctx, customoverlaysrepository.CreateInput{
			ChannelID: input.ChannelID,
			Name:      input.Name,
			Width:     input.Width,
			Height:    input.Height,
			Layers:    layers,
		},
	)
	if err != nil {
		return entity.CustomOverlay{}, err
	}

	return c.mapToEntity(overlay), nil
}

type UpdateInput struct {
	ChannelID string
	Name      string
	Width     int
	Height    int
	Layers    []LayerInput
}

func (c *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateInput,
) (entity.CustomOverlay, error) {
	layers := make([]customoverlaysrepository.LayerInput, len(input.Layers))
	for i, layer := range input.Layers {
		layers[i] = customoverlaysrepository.LayerInput{
			Type:                    layer.Type,
			Width:                   layer.Width,
			Height:                  layer.Height,
			CreatedAt:               layer.CreatedAt,
			UpdatedAt:               layer.UpdatedAt,
			PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
			TransformString:         layer.TransformString,
			SettingsHtmlHtml:        layer.SettingsHtmlHtml,
			SettingsHtmlCss:         layer.SettingsHtmlCss,
			SettingsHtmlJs:          layer.SettingsHtmlJs,
			SettingsHtmlDataPoll:    layer.SettingsHtmlDataPoll,
			SettingsImageURL:        layer.SettingsImageURL,
		}
	}

	overlay, err := c.repository.Update(
		ctx, id, customoverlaysrepository.UpdateInput{
			ChannelID: input.ChannelID,
			Name:      input.Name,
			Width:     input.Width,
			Height:    input.Height,
			Layers:    layers,
		},
	)
	if err != nil {
		return entity.CustomOverlay{}, err
	}

	overlayID := id.String()

	_, err = c.websocketsGrpc.RefreshOverlaySettings(
		ctx,
		&websockets.RefreshOverlaysRequest{
			ChannelId: input.ChannelID,
			OverlayId: &overlayID,
		},
	)
	if err != nil {
		return entity.CustomOverlay{}, err
	}

	return c.mapToEntity(overlay), nil
}

func (c *Service) Delete(ctx context.Context, id uuid.UUID, dashboardID string) error {
	overlay, err := c.repository.GetOne(ctx, id)
	if err != nil {
		return err
	}

	if overlay.ChannelID != dashboardID {
		return fmt.Errorf("overlay with id %s does not belong to dashboard with id %s", id, dashboardID)
	}

	return c.repository.Delete(ctx, id)
}

func (c *Service) ParseTextInHtml(ctx context.Context, dashboardID, text string) (string, error) {
	res, err := c.twirBus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID: dashboardID,
			Text:      text,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Data.Text, nil
}
