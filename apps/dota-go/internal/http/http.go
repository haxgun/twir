package http

import (
	"context"
	"net/http"

	"github.com/satont/twir/apps/dota-go/internal/processor"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type Http struct {
	logger    *logrus.Logger
	processor *processor.Processor
}

type Opts struct {
	fx.In
	Lc        fx.Lifecycle
	Logger    *logrus.Logger
	Processor *processor.Processor
}

func New(opts Opts) {
	service := &Http{
		logger:    opts.Logger,
		processor: opts.Processor,
	}

	http.Handle(
		"/", service,
	)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					opts.Logger.Infoln("Http server started and listening on http://0.0.0.0:3000")
					err := http.ListenAndServe("0.0.0.0:3008", nil)
					if err != nil {
						opts.Logger.Fatal(err)
					}
				}()
				return nil
			},
			OnStop: nil,
		},
	)
}
