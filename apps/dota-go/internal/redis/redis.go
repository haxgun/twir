package redis

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/sirupsen/logrus"
)

func New(config cfg.Config, logger *logrus.Logger) (*redis.Client, error) {
	params, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(params)
	logger.Infoln("Redis connected")

	return client, nil
}
