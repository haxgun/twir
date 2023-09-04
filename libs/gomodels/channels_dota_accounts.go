package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type ChannelsDotaAccounts struct {
	ID        string `gorm:"primary_key;column:id;type:TEXT;"        json:"id"`
	ChannelID string `gorm:"primary_key;column:channelId;type:TEXT;" json:"channel_id"`
	NickName  string `gorm:"column:nickname;type:TEXT;"              json:"nick_name"`
}

func (c ChannelsDotaAccounts) TableName() string {
	return "channels_dota_accounts"
}
