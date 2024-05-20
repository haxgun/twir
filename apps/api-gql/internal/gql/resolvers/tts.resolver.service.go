package resolvers

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

type RHVoiceInfoResponse struct {
	DEFAULTFORMAT            string            `json:"DEFAULT_FORMAT"`
	DEFAULTVOICE             string            `json:"DEFAULT_VOICE"`
	FORMATS                  map[string]string `json:"FORMATS"`
	SUPPORTVOICES            []string          `json:"SUPPORT_VOICES"`
	RhvoiceWrapperApiVersion string            `json:"rhvoice_wrapper_api_version"`
	RhvoiceWrapperCmd        struct {
		Flac []string `json:"flac"`
		Mp3  []string `json:"mp3"`
		Opus []string `json:"opus"`
	} `json:"rhvoice_wrapper_cmd"`
	RhvoiceWrapperLibraryVersion string                 `json:"rhvoice_wrapper_library_version"`
	RhvoiceWrapperProcess        bool                   `json:"rhvoice_wrapper_process"`
	RhvoiceWrapperThreadCount    int                    `json:"rhvoice_wrapper_thread_count"`
	RhvoiceWrapperVoiceProfiles  []string               `json:"rhvoice_wrapper_voice_profiles"`
	RhvoiceWrapperVoicesInfo     map[string]RHVoiceInfo `json:"rhvoice_wrapper_voices_info"`
}

type RHVoiceInfo struct {
	Country string `json:"country"`
	Gender  string `json:"gender"`
	Lang    string `json:"lang"`
	Name    string `json:"name"`
	No      int    `json:"no"`
}

func (r *tTSVoicesResolver) ttsRhVoiceGetInfo(ctx context.Context) (*RHVoiceInfoResponse, error) {
	var data *RHVoiceInfoResponse
	resp, err := req.
		R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get(fmt.Sprintf("http://%s/info", r.config.TTSServiceUrl))
	if err != nil {
		return nil, fmt.Errorf("tts service is not available: %w", err)
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("tts service is not available: %w", err)
	}

	return data, nil
}

func ttsVoiceDbServiceToGql(service model.ChannelTTSVoiceService) gqlmodel.TTSSettingsVoiceService {
	switch service {
	case model.ChannelTTSVoiceServiceRH:
		return gqlmodel.TTSSettingsVoiceServiceRh
	case model.ChannelTTSVoiceServiceSilero:
		return gqlmodel.TTSSettingsVoiceServiceSilero
	}
	return gqlmodel.TTSSettingsVoiceServiceRh
}

func ttsVoiceGqlServiceToDb(service gqlmodel.TTSSettingsVoiceService) model.ChannelTTSVoiceService {
	switch service {
	case gqlmodel.TTSSettingsVoiceServiceRh:
		return model.ChannelTTSVoiceServiceRH
	case gqlmodel.TTSSettingsVoiceServiceSilero:
		return model.ChannelTTSVoiceServiceSilero
	}
	return model.ChannelTTSVoiceServiceRH
}

func (r *queryResolver) ttsGetChannelSettings(
	ctx context.Context,
) (*gqlmodel.TTSSettings, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entity model.ChannelsTTS
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		Preload("DisallowedVoices").
		First(&entity).Error; err != nil {
		return nil, fmt.Errorf("failed to get tts settings: %w", err)
	}

	disallowedVoices := make([]gqlmodel.TTSSettingsDisallowedVoice, 0, len(entity.DisallowedVoices))
	for _, voice := range entity.DisallowedVoices {
		disallowedVoices = append(
			disallowedVoices,
			gqlmodel.TTSSettingsDisallowedVoice{
				Voice:   voice.Voice,
				Service: ttsVoiceDbServiceToGql(voice.Service),
			},
		)
	}

	return &gqlmodel.TTSSettings{
		Enabled:                            entity.Enabled,
		Rate:                               entity.Rate,
		Volume:                             entity.Volume,
		Pitch:                              entity.Pitch,
		Voice:                              entity.Voice,
		VoiceService:                       ttsVoiceDbServiceToGql(entity.VoiceService),
		AllowUsersChooseVoiceInMainCommand: entity.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         entity.MaxSymbols,
		DisallowedVoices:                   disallowedVoices,
		DoNotReadEmoji:                     entity.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              entity.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     entity.DoNotReadLinks,
		ReadChatMessages:                   entity.ReadChatMessages,
		ReadChatMessagesNicknames:          entity.ReadChatMessagesNicknames,
	}, nil
}

func (r *tTSSettingsResolver) getUsersSettings(
	ctx context.Context,
) ([]gqlmodel.TTSUserSettings, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelTTSUserSettings
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to get tts settings: %w", err)
	}

	result := make([]gqlmodel.TTSUserSettings, 0, len(entities))

	for _, entity := range entities {
		result = append(
			result,
			gqlmodel.TTSUserSettings{
				UserID:       entity.UserID,
				Voice:        entity.Voice,
				VoiceService: ttsVoiceDbServiceToGql(entity.Service),
				Rate:         entity.Rate,
				Pitch:        entity.Pitch,
			},
		)
	}

	return result, nil
}
