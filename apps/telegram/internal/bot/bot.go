package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kr/pretty"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/telegram/internal/bot/middlewares"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config      config.Config
	Logger      logger.Logger
	Middlewares *middlewares.Middlewares
}

func New(opts Opts) (*Bot, error) {
	t := &Bot{
		config:      opts.Config,
		middlewares: opts.Middlewares,
	}

	if opts.Config.TelegramBotApiToken == "" {
		opts.Logger.Warn("Telegram bot token is empty, bot will not be started")
		return t, nil
	}

	if err := t.createBot(); err != nil {
		return nil, err
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := t.Start(); err != nil {
					return err
				}

				me, err := t.bot.GetMe(context.Background())
				if err != nil {
					return err
				}

				opts.Logger.Info("Bot started", slog.String("username", me.Username))

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return t.Stop()
			},
		},
	)

	return t, nil
}

type Bot struct {
	config      config.Config
	middlewares *middlewares.Middlewares

	stopFunc func()
	botCtx   context.Context
	bot      *bot.Bot
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	link, _ := b.CreateInvoiceLink(
		ctx, &bot.CreateInvoiceLinkParams{
			Title:         "Test",
			Description:   "Test desc",
			Payload:       "q",
			ProviderToken: "",
			Currency:      "XTR",
			Prices: []models.LabeledPrice{
				{
					Label:  "Test",
					Amount: 1,
				},
			},
			SubscriptionPeriod:        2592000, // 30 days
			MaxTipAmount:              0,
			SuggestedTipAmounts:       nil,
			ProviderData:              "",
			PhotoURL:                  "",
			PhotoSize:                 0,
			PhotoWidth:                0,
			PhotoHeight:               0,
			NeedName:                  false,
			NeedPhoneNumber:           false,
			NeedEmail:                 false,
			NeedShippingAddress:       false,
			SendPhoneNumberToProvider: false,
			SendEmailToProvider:       false,
			IsFlexible:                false,
		},
	)

	b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   link,
		},
	)
}

func (c *Bot) createBot() error {
	botCtx, cancel := context.WithCancel(context.Background())
	c.botCtx = botCtx
	c.stopFunc = cancel

	opts := []bot.Option{
		bot.WithMiddlewares(
			c.middlewares.LogIncomingUpdate,
			c.middlewares.CheckIsLoggedIn,
		),
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(c.config.TelegramBotApiToken, opts...)
	if err != nil {
		return err
	}

	const chargeId = "stxfig_4W2t40p698FgpVSMAdkKNR16jgWs-KEtSmODXWR-wvkd3l-zziOU6LwbfsBZQxDo8GjClD5u9YnzAauCYVukQalTgzYVgFYeEl5KRgIKCb0AQKB_WSHx0KrRL6rj"

	go func() {
		fmt.Println(
			b.RefundStarPayment(
				context.Background(), &bot.RefundStarPaymentParams{
					UserID:                  210144787,
					TelegramPaymentChargeID: chargeId,
				},
			),
		)

		fmt.Println(
			b.EditUserStarSubscription(
				context.Background(),
				&bot.EditUserStarSubscriptionParams{
					UserID:                  210144787,
					TelegramPaymentChargeID: chargeId,
					IsCanceled:              true,
				},
			),
		)
	}()

	b.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update.PreCheckoutQuery != nil
		}, func(ctx context.Context, q *bot.Bot, update *models.Update) {
			q.AnswerPreCheckoutQuery(
				ctx, &bot.AnswerPreCheckoutQueryParams{
					PreCheckoutQueryID: update.PreCheckoutQuery.ID,
					OK:                 true,
				},
			)
		},
	)

	b.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update.Message.RefundedPayment != nil
		}, func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			pretty.Println(update.Message.RefundedPayment)
		},
	)

	b.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update.Message.SuccessfulPayment != nil
		},
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			pretty.Println(update.Message.SuccessfulPayment)
		},
	)

	c.bot = b
	return nil
}

func (c *Bot) Start() error {
	go c.bot.Start(c.botCtx)
	return nil
}

func (c *Bot) Stop() error {
	c.stopFunc()
	return nil
}

func (c *Bot) CreateInvoiceLink(ctx context.Context, title string, price int) string {
	link, _ := c.bot.CreateInvoiceLink(
		ctx,
		&bot.CreateInvoiceLinkParams{
			Title:         title,
			Description:   title,
			Payload:       title,
			ProviderToken: "",
			Currency:      "XTR",
			Prices: []models.LabeledPrice{
				{
					Label:  "Test",
					Amount: price,
				},
			},
			SubscriptionPeriod:        2592000, // 30 days
			MaxTipAmount:              0,
			SuggestedTipAmounts:       nil,
			ProviderData:              "",
			PhotoURL:                  "",
			PhotoSize:                 0,
			PhotoWidth:                0,
			PhotoHeight:               0,
			NeedName:                  false,
			NeedPhoneNumber:           false,
			NeedEmail:                 false,
			NeedShippingAddress:       false,
			SendPhoneNumberToProvider: false,
			SendEmailToProvider:       false,
			IsFlexible:                false,
		},
	)

	return link
}
