package tts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type settings struct {
	Enabled      bool
	Rate         int
	Volume       int
	Pitch        int
	Voice        string
	VoiceService model.ChannelTTSVoiceService
}

func getSettings(
	ctx context.Context,
	db *gorm.DB,
	channelId string,
	userId *string,
) *settings {
	channelEntity := &model.ChannelsTTS{}
	userEntity := &model.ChannelsTTSUserSettings{}

	if userId != nil {
		err := db.WithContext(ctx).Where(
			"channel_id = ? AND user_id = ?",
			channelId,
			userId,
		).Find(userEntity).Error
		if err != nil {
			return nil
		}
	} else {
		err := db.WithContext(ctx).Where("channel_id = ?", channelId).Find(channelEntity).Error
		if err != nil {
			return nil
		}
	}

	if userEntity.UserID != "" {
		return &settings{
			Enabled:      userEntity.UserID != "",
			Rate:         userEntity.Rate,
			Volume:       userEntity.Volume,
			Pitch:        userEntity.Pitch,
			Voice:        userEntity.Voice,
			VoiceService: userEntity.Service,
		}
	} else if channelEntity.ChannelID != "" {
		return &settings{
			Enabled:      channelEntity.Enabled,
			Rate:         channelEntity.Rate,
			Volume:       channelEntity.Volume,
			Pitch:        channelEntity.Pitch,
			Voice:        channelEntity.Voice,
			VoiceService: channelEntity.VoiceService,
		}
	} else {
		return nil
	}
}

type Voice struct {
	Name    string
	Country string
}

func getVoices(ctx context.Context, cfg *config.Config) []Voice {
	data := map[string]any{}
	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get(fmt.Sprintf("http://%s/info", cfg.TTSServiceUrl))
	if err != nil {
		return nil
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		return nil
	}

	parsedJson := gjson.ParseBytes(bytes)
	voices := []Voice{}
	parsedJson.Get("rhvoice_wrapper_voices_info").ForEach(
		func(key, value gjson.Result) bool {
			voices = append(
				voices, Voice{
					Name:    key.String(),
					Country: value.Get("country").String(),
				},
			)

			return true
		},
	)

	return voices
}

func updateSettings(
	ctx context.Context,
	db *gorm.DB,
	entity *model.ChannelModulesSettings,
	settings *settings,
) error {
	bytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return db.
		Model(entity).
		WithContext(ctx).
		Updates(map[string]interface{}{"settings": bytes}).
		Error
}

func createUserSettings(
	ctx context.Context,
	db *gorm.DB,
	rate,
	pitch int,
	voice,
	channelId,
	userId string,
) (
	*model.ChannelModulesSettings,
	*modules.TTSSettings,
	error,
) {
	userModel := &model.ChannelModulesSettings{
		ID:        uuid.New().String(),
		Type:      "tts",
		Settings:  nil,
		ChannelId: channelId,
		UserId:    null.StringFrom(userId),
	}

	userSettings := &modules.TTSSettings{
		Enabled: lo.ToPtr(true),
		Rate:    rate,
		Volume:  70,
		Pitch:   pitch,
		Voice:   voice,
	}

	bytes, err := json.Marshal(userSettings)
	if err != nil {
		return nil, nil, err
	}

	userModel.Settings = bytes

	err = db.WithContext(ctx).Create(userModel).Error
	if err != nil {
		return nil, nil, err
	}

	return userModel, userSettings, nil
}

func switchEnableState(ctx context.Context, db *gorm.DB, channelId string, newState bool) error {
	channelSettings, channelModele := getSettings(ctx, db, channelId, "")

	if channelSettings == nil {
		return errors.New("tts not configured")
	}

	channelSettings.Enabled = &newState
	err := updateSettings(ctx, db, channelModele, channelSettings)
	if err != nil {
		return err
	}

	return nil
}

func isValidUrl(input string) bool {
	if u, e := url.Parse(input); e == nil {
		if u.Host != "" {
			return dnsCheck(u.Host)
		}

		return dnsCheck(input)
	}

	return false
}

func dnsCheck(input string) bool {
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")

	ips, _ := net.LookupIP(input)
	return len(ips) > 0
}
