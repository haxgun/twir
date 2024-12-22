package streams

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/redis_keys"
	"github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorymodel "github.com/twirapp/twir/libs/repositories/streams/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Redis             *redis.Client
	StreamsRepository streams.Repository
}

func New(opts Opts) *Service {
	return &Service{
		redis:             opts.Redis,
		streamsRepository: opts.StreamsRepository,
	}
}

type Service struct {
	redis             *redis.Client
	streamsRepository streams.Repository
}

func (c *Service) mapToEntity(m streamsrepositorymodel.Stream) entity.Stream {
	return entity.Stream{
		ID:           m.ID,
		UserID:       m.UserID,
		UserLogin:    m.UserLogin,
		UserName:     m.UserName,
		GameID:       m.GameID,
		GameName:     m.GameName,
		CommunityIDs: m.CommunityIDs,
		Type:         m.Type,
		Title:        m.Title,
		ViewerCount:  m.ViewerCount,
		StartedAt:    m.StartedAt,
		Language:     m.Language,
		ThumbnailUrl: m.ThumbnailUrl,
		TagIDs:       m.TagIDs,
		Tags:         m.Tags,
		IsMature:     m.IsMature,
	}
}

func (c *Service) GetByChannelID(ctx context.Context, channelID string) (entity.Stream, error) {
	stream, err := c.streamsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return entity.StreamNil, err
	}

	return c.mapToEntity(stream), nil
}

func (c *Service) GetParsedChatMessages(ctx context.Context, streamID string) (int, error) {
	parsedMessages, err := c.redis.Get(
		ctx,
		redis_keys.StreamParsedMessages(
			streamID,
		),
	).Int()
	if err != nil {
		return 0, fmt.Errorf("cannot get parsed chat messages: %w", err)
	}

	return parsedMessages, nil
}
