package overlays

import (
	"github.com/satont/twir/apps/api/internal/impl_deps"
)

type Overlays struct {
	*impl_deps.Deps
}

// func (c *Overlays) OverlaysParseHtml(ctx context.Context, req *overlays.ParseHtmlOverlayRequest) (
// 	*overlays.ParseHtmlOverlayResponse, error,
// ) {
// 	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	res, err := c.Bus.Parser.ParseVariablesInText.Request(
// 		ctx, parser.ParseVariablesInTextRequest{
// 			ChannelID: dashboardId,
// 			Text:      base64ToText(req.GetHtml()),
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &overlays.ParseHtmlOverlayResponse{
// 		Html: textToBase64(res.Data.Text),
// 	}, nil
// }
