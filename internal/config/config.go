package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultModelID = "scribe_v2"

type Config struct {
	TelegramToken    string
	ElevenLabsAPIKey string
	ElevenLabsModel  string
	allowedChats     map[int64]struct{}
}

func Load() (*Config, error) {
	cfg := &Config{
		TelegramToken:    strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		ElevenLabsAPIKey: strings.TrimSpace(os.Getenv("ELEVENLABS_API_KEY")),
		ElevenLabsModel:  strings.TrimSpace(os.Getenv("ELEVENLABS_MODEL_ID")),
		allowedChats:     make(map[int64]struct{}),
	}

	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.ElevenLabsAPIKey == "" {
		return nil, fmt.Errorf("ELEVENLABS_API_KEY is required")
	}
	if cfg.ElevenLabsModel == "" {
		cfg.ElevenLabsModel = defaultModelID
	}

	raw := os.Getenv("ALLOWED_CHAT_ID")
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ALLOWED_CHAT_ID: %q is not a chat id", part)
		}
		cfg.allowedChats[id] = struct{}{}
	}

	return cfg, nil
}

func (c *Config) ChatAllowed(chatID int64) bool {
	if len(c.allowedChats) == 0 {
		return true
	}
	_, ok := c.allowedChats[chatID]
	return ok
}
