package steam

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/paralin/go-steam"
	"github.com/paralin/go-steam/protocol/steamlang"
	cfg "github.com/satont/twir/libs/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type Client struct {
	client *steam.Client
	ready  chan bool
}

func New(lc fx.Lifecycle, config cfg.Config, logger *logrus.Logger) *Client {
	localLog := logger.WithField("username", config.SteamUserName)

	client := &Client{
		client: steam.NewClient(),
		ready:  make(chan bool),
	}

	if err := steam.InitializeSteamDirectory(); err != nil {
		panic(err)
	}

	lastReconnect := time.Time{}
	reconnectIssues := 0

	go func() {
		ch := client.client.Events()
		for event := range ch {
			switch event.(type) {
			case *steam.ConnectedEvent:
				localLog.Infoln("Steam: connected, logging")
				client.client.Auth.LogOn(
					&steam.LogOnDetails{
						Username: config.SteamUserName,
						Password: config.SteamPassword,
					},
				)
			case *steam.DisconnectedEvent:
				localLog.Infoln("Steam: disconnected")
			case *steam.MachineAuthUpdateEvent:
			case *steam.LogOnFailedEvent:
				localLog.WithError(fmt.Errorf(spew.Sdump(event))).Fatalf("steam login failed")
			case *steam.LoggedOnEvent:
				client.ready <- true
				localLog.Infoln("Steam: logged in")
				client.client.Social.SetPersonaState(steamlang.EPersonaState_Online)
				if time.Since(lastReconnect) < 30*time.Second {
					reconnectIssues++
				} else {
					reconnectIssues = 0
				}

				if reconnectIssues > 5 {
					localLog.Errorln("Steam: reconnect issues detected, restarting bot")
				}
			}
		}
	}()

	go func() {
		tick := time.NewTicker(time.Second * 10)
		failedCount := 0
		for {
			select {
			case <-tick.C:
				if !client.client.Connected() {
					localLog.Warningln("Steam: reconnecting")
					client.client.Connect()
					failedCount++
					if failedCount > 10 {
						localLog.Errorln("Steam: failed to connect after 10 attempts")
					}
				} else {
					failedCount = 0
				}
			}
		}
	}()

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				if client.client.Connected() {
					client.client.Disconnect()
				}
				return nil
			},
			OnStart: func(ctx context.Context) error {
				localLog.Infoln("Steam: connecting")
				client.client.Connect()
				return nil
			},
		},
	)

	return client
}

func (c *Client) Raw() *steam.Client {
	return c.client
}

func (c *Client) Ready() <-chan bool {
	return c.ready
}

func (c *Client) IsConnected() bool {
	return c.client.Connected()
}
