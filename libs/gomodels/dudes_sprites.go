package model

import (
	"github.com/google/uuid"
)

type DudeSprite struct {
	ID        uuid.UUID
	Name      string
	Published bool
	OwnerID   string

	Owner *Users
}

type DudeSpriteLayerType string

const (
	DudeSpriteLayerTypeBody      DudeSpriteLayerType = "body"
	DudeSpriteLayerTypeEyes      DudeSpriteLayerType = "eyes"
	DudeSpriteLayerTypeMouth     DudeSpriteLayerType = "mouth"
	DudeSpriteLayerTypeCosmetics DudeSpriteLayerType = "cosmetics"
)

type DudeSpriteLayer struct {
	ID    uuid.UUID
	Type  DudeSpriteLayerType
	Name  string
	Color string
	Image []byte

	Approved   bool
	ApprovedBy string

	ApprovedByUser *Users
}

type DudeUserChannelSprite struct {
	ID        uuid.UUID
	UserID    string
	ChannelID string
	SpriteID  uuid.UUID

	Approved       bool
	ApprovedBy     string
	ApprovedByUser *Users

	// relations
	User    *Users
	Channel *Channels
	Sprite  *DudeSprite
}
