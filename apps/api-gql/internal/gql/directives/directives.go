package directives

import (
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Sessions *sessions.Sessions
	Gorm     *gorm.DB
	Redis    *redis.Client
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions: opts.Sessions,
		gorm:     opts.Gorm,
		redis:    opts.Redis,
	}
}

type Directives struct {
	sessions *sessions.Sessions
	gorm     *gorm.DB
	redis    *redis.Client
}
