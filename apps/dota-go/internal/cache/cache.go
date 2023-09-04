package cache

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const prefix = "dota2:channels:"

type Cache struct {
	redis *redis.Client
	db    *gorm.DB
}

type Opts struct {
	fx.In

	Redis *redis.Client
	DB    *gorm.DB
}

func New(opts Opts) *Cache {
	return &Cache{
		redis: opts.Redis,
		db:    opts.DB,
	}
}
