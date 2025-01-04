package middlewares

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Middlewares) CheckIsLoggedIn(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		if update.Message == nil {
			c.logger.Debug("CheckIsLoggedIn: update.Message is nil, skipping")
			next(ctx, bot, update)
			return
		}

		isLoggedIn, err := c.telegramService.CheckIsIntegrationExistsByChatID(
			ctx,
			update.Message.Chat.ID,
		)
		if err != nil {
			c.logger.Error("CheckIsIntegrationExistsByChatID", slog.Any("err", err))
			return
		}

		if !isLoggedIn {
			c.logger.Warn("User is not logged in", slog.Any("chatID", update.Message.Chat))
			return
		}

		next(ctx, bot, update)
	}
}
