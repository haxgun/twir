package model

import (
	"github.com/google/uuid"
)

type ChannelTTSVoiceService string

const (
	ChannelTTSVoiceServiceRH     ChannelTTSVoiceService = "rhvoice"
	ChannelTTSVoiceServiceSilero ChannelTTSVoiceService = "silero"
)

type ChannelsTTS struct {
	ID                                 uuid.UUID                    `gorm:"column:id;primary_key;type:uuid"`
	ChannelID                          string                       `gorm:"column:channel_id;type:text"`
	Enabled                            bool                         `gorm:"column:enabled;type:bool"`
	Rate                               int                          `gorm:"column:rate;type:int"`
	Volume                             int                          `gorm:"column:volume;type:int"`
	Pitch                              int                          `gorm:"column:pitch;type:int"`
	Voice                              string                       `gorm:"column:voice;type:text"`
	VoiceService                       ChannelTTSVoiceService       `gorm:"column:voice_service;type:text"`
	AllowUsersChooseVoiceInMainCommand bool                         `gorm:"column:allow_users_choose_voice_in_main_command;type:bool"`
	MaxSymbols                         int                          `gorm:"column:max_symbols;type:int"`
	DoNotReadEmoji                     bool                         `gorm:"column:do_not_read_emoji;type:bool"`
	DoNotReadTwitchEmotes              bool                         `gorm:"column:do_not_read_twitch_emotes;type:bool"`
	DoNotReadLinks                     bool                         `gorm:"column:do_not_read_links;type:bool"`
	ReadChatMessages                   bool                         `gorm:"column:read_chat_messages;type:bool"`
	ReadChatMessagesNicknames          bool                         `gorm:"column:read_chat_messages_nicknames;type:bool"`
	DisallowedVoices                   []ChannelsTTSDisallowedVoice `gorm:"foreignKey:ChannelID"`
}

func (ChannelsTTS) TableName() string {
	return "channels_overlays_tts"
}

type ChannelsTTSDisallowedVoice struct {
	ID        uuid.UUID              `gorm:"column:id;primary_key;type:uuid"`
	ChannelID string                 `gorm:"column:channel_id;type:text"`
	Voice     string                 `gorm:"column:voice;type:text"`
	Service   ChannelTTSVoiceService `gorm:"column:service;type:text"`
}

func (ChannelsTTSDisallowedVoice) TableName() string {
	return "channels_overlays_tts_disallowed_voices"
}

type ChannelTTSUserSettings struct {
	ID        uuid.UUID              `gorm:"column:id;primary_key;type:uuid"`
	ChannelID string                 `gorm:"column:channel_id;type:text"`
	UserID    string                 `gorm:"column:user_id;type:text"`
	Voice     string                 `gorm:"column:voice;type:text"`
	Service   ChannelTTSVoiceService `gorm:"column:service;type:text"`
	Rate      int                    `gorm:"column:rate;type:int"`
	Pitch     int                    `gorm:"column:pitch;type:int"`
	Volume    int                    `gorm:"column:volume;type:int"`
}

func (ChannelTTSUserSettings) TableName() string {
	return "channels_overlays_tts_users_settings"
}
