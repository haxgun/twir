package entity

import (
	"time"
)

type DashboardStats struct {
	CategoryID     string
	CategoryName   string
	Viewers        *int
	StartedAt      *time.Time
	Title          string
	ChatMessages   int
	Followers      int
	UsedEmotes     int
	RequestedSongs int
	Subs           int
}

var DashboardStatsNil = DashboardStats{}
