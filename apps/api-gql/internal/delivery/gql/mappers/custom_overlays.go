package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func CustomOverlayToGql(e entity.CustomOverlay) gqlmodel.CustomOverlay {
	layers := make([]gqlmodel.CustomOverlayLayer, len(e.Layers))
	for i, layer := range e.Layers {
		layers[i] = gqlmodel.CustomOverlayLayer{
			ID:                      layer.ID,
			Type:                    gqlmodel.CustomOverlayLayerType(layer.Type),
			OverlayID:               layer.OverlayID,
			Width:                   layer.Width,
			Height:                  layer.Height,
			CreatedAt:               layer.CreatedAt,
			UpdatedAt:               layer.UpdatedAt,
			PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
			TransformString:         layer.TransformString,
		}

		if layer.SettingsHtmlHtml != nil && layer.SettingsHtmlCss != nil && layer.SettingsHtmlJs != nil {
			layers[i].Settings = gqlmodel.CustomOverlayLayerSettingsHTML{
				EmptyHolder:         nil,
				HTML:                *layer.SettingsHtmlHtml,
				CSS:                 *layer.SettingsHtmlCss,
				Js:                  *layer.SettingsHtmlJs,
				PollSecondsInterval: *layer.SettingsHtmlDataPoll,
			}
		}

		if layer.SettingsImageURL != nil {
			layers[i].Settings = gqlmodel.CustomOverlayLayerSettingsImage{
				EmptyHolder: nil,
				URL:         *layer.SettingsImageURL,
			}
		}
	}

	return gqlmodel.CustomOverlay{
		ID:        e.ID,
		ChannelID: e.ChannelID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		Width:     e.Width,
		Height:    e.Height,
		Layers:    layers,
	}
}
