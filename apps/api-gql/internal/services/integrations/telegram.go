package integrations

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	channelsintegrationstelegramrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_telegram"
)

func (c *Service) TelegramVerifyStringSignature(inputQuery url.Values) bool {
	pairs := make([]string, 0, len(inputQuery))
	var hash string

	for k, v := range inputQuery {
		// Store found sign.
		if k == "hash" {
			hash = v[0]
			continue
		}

		// Append new pair.
		pairs = append(pairs, k+"="+v[0])
	}
	sort.Strings(pairs)

	keyHash := sha256.New()
	keyHash.Write([]byte(c.config.TelegramBotApiToken))
	secretKey := keyHash.Sum(nil)

	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(strings.Join(pairs, "\n")))
	calculatedHash := hex.EncodeToString(h.Sum(nil))
	fmt.Println(pairs)

	if calculatedHash != hash {
		return false
	}

	return true
}

type TelegramUserIntegrationInput struct {
	ChannelID      string
	TelegramUserID string
}

func (c *Service) TelegramCreateUserIntegration(
	ctx context.Context,
	input TelegramUserIntegrationInput,
) error {
	_, err := c.telegramRepository.Create(
		ctx, channelsintegrationstelegramrepository.CreateInput{
			ChannelID:      input.ChannelID,
			TelegramChatID: input.TelegramUserID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create telegram integration: %w", err)
	}

	return nil
}
