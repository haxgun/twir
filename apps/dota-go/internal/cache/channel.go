package cache

import (
	"context"
	"time"

	model "github.com/satont/twir/libs/gomodels"
)

type Channel struct {
	ID string
}

func (c *Cache) ChannelByToken(ctx context.Context, apiKey string) (Channel, error) {
	channel := Channel{}

	isExists, _ := c.redis.Exists(ctx, prefix+apiKey).Result()
	if isExists == 0 {
		user := model.Users{}
		if err := c.db.Where(`"apiKey" = ?`, apiKey).First(&user).Error; err != nil {
			return channel, err
		}
		channel.ID = user.ID
		if err := c.redis.Set(ctx, prefix+apiKey, user.ID, 7*24*time.Hour).Err(); err != nil {
			return channel, err
		}
	}

	res, err := c.redis.Get(ctx, prefix+apiKey).Result()
	if err != nil {
		return channel, err
	}
	channel.ID = res
	return channel, nil
}
