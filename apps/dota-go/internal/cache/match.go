package cache

import (
	"context"
	"fmt"
)

func buildMatchKey(channelID, accountID, matchID string) string {
	return fmt.Sprintf("%s:%s:accounts:%s:match:%s", prefix, channelID, accountID, matchID)
}

func (c *Cache) MatchIsExists(ctx context.Context, channelID, accountID, matchID string) (
	bool,
	error,
) {
	exists, err := c.redis.Exists(
		ctx,
		buildMatchKey(channelID, accountID, matchID),
	).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (c *Cache) MatchGetAll(ctx context.Context, channelID, accountID string) ([]string, error) {
	var keys []string
	// scan all keys
	iter := c.redis.Scan(ctx, 0, buildMatchKey(channelID, accountID, "*"), 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
