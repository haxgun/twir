package main

import (
	"github.com/satont/twir/apps/dota-go/internal/cache"
	"github.com/satont/twir/apps/dota-go/internal/dota"
	"github.com/satont/twir/apps/dota-go/internal/gorm"
	"github.com/satont/twir/apps/dota-go/internal/http"
	"github.com/satont/twir/apps/dota-go/internal/processor"
	"github.com/satont/twir/apps/dota-go/internal/redis"
	"github.com/satont/twir/apps/dota-go/internal/steam"
	cfg "github.com/satont/twir/libs/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			// author of go-steam and go-dota2 - dolboeb
			func() *logrus.Logger {
				return logrus.New()
			},
			cfg.NewFx,
			gorm.New,
			redis.New,
			processor.New,
			cache.New,
			steam.New,
			dota.New,
		),
		fx.Invoke(
			dota.New,
			http.New,
		),
	).Run()
}
