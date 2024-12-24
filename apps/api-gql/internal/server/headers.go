package server

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
)

func SetHeader(ctx context.Context, key, value string) error {
	gin, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return err
	}

	gin.Header(key, value)

	return nil
}
