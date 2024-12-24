package ttsvoices

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/imroc/req/v3"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config config.Config
}

func New(opts Opts) *Service {
	s := &Service{
		config: opts.Config,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := s.fetchRHVoices(ctx); err != nil {
					return err
				}

				return nil
			},
		},
	)

	return s
}

type Service struct {
	config config.Config

	rhVoices []entity.TTSRHVoice
}

type rhvoicesResponse struct {
	Voices map[string]rhvoicesResponseVoice `json:"rhvoice_wrapper_voices_info"`
}

type rhvoicesResponseVoice struct {
	Country string `json:"country"`
	Gender  string `json:"gender"`
	Lang    string `json:"lang"`
	Name    string `json:"name"`
	No      int    `json:"no"`
}

func (c *Service) fetchRHVoices(ctx context.Context) error {
	var result rhvoicesResponse
	resp, err := req.
		R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://%s/info", c.config.TTSServiceUrl))
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("cannot get rh voices: %s", resp.String())
	}

	for key, value := range result.Voices {
		c.rhVoices = append(
			c.rhVoices,
			entity.TTSRHVoice{
				Code:    key,
				Country: value.Country,
				Gender:  value.Gender,
				Lang:    value.Lang,
				Name:    value.Name,
				No:      value.No,
			},
		)
	}

	slices.SortFunc(
		c.rhVoices,
		func(a, b entity.TTSRHVoice) int {
			return strings.Compare(a.Country, b.Country)
		},
	)

	return nil
}

func (c *Service) GetRHVoices() []entity.TTSRHVoice {
	return c.rhVoices
}
