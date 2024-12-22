package entity

import (
	"time"
)

type Stream struct {
	ID           string
	UserID       string
	UserLogin    string
	UserName     string
	GameID       string
	GameName     string
	CommunityIDs []string
	Type         string
	Title        string
	ViewerCount  int
	StartedAt    time.Time
	Language     string
	ThumbnailUrl string
	TagIDs       []string
	Tags         []string
	IsMature     bool
}

var StreamNil = Stream{}
