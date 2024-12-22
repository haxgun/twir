package streams

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/streams/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.Stream, error)
}
