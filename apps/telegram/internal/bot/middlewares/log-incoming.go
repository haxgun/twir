package middlewares

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Middlewares) LogIncomingUpdate(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		var slogAttrs []slog.Attr
		if update.Message != nil {
			slogAttrs = append(
				slogAttrs,
				slog.Group(
					"user",
					slog.String("username", update.Message.Chat.Username),
					slog.Int64("id", update.Message.Chat.ID),
				),
			)
		}

		slogAttrs = append(slogAttrs, slog.Int64("id", update.ID))

		c.logger.Info(
			"Incoming update",
			slogAttrs,
		)

		next(ctx, bot, update)
	}
}
