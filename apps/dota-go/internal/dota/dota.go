package dota

import (
	"context"
	"time"

	"github.com/paralin/go-dota2"
	"github.com/satont/twir/apps/dota-go/internal/steam"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle, steam *steam.Client, logger *logrus.Logger) *dota2.Dota2 {
	shutUpLogger := logrus.New()
	shutUpLogger.SetLevel(logrus.FatalLevel)

	client := dota2.New(steam.Raw(), shutUpLogger)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					for {
						<-steam.Ready()
						client.SetPlaying(true)
						time.Sleep(1 * time.Second)
						client.SayHello()
						logger.Infoln("Dota: opened and connected")
					}
				}()
				return nil
			},
			OnStop: nil,
		},
	)

	go func() {
		for {
			if steam.IsConnected() {
				client.SayHello()
			}

			time.Sleep(2 * time.Second)
		}
	}()

	return client
}
