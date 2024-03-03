package config

import (
	"errors"
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
)

type Config struct {
	TGBotConf struct {
		Token             string
		UpdatesBufferSize int
		PermittedChatID   int
		PollIntervalSec   int
	}
	LinkdingConf struct {
		LinkdingAddr string
		UserToken    string
	}
	LogLevel string
	LogFile  string
}

func GetConfig() (*Config, error) {
	k := koanf.New(".")

	defaults := map[string]interface{}{
		"LGTG_BOT_TOKEN":                "",
		"LGTG_BOT_UPDATES_BUFFER_SIZE":  1,
		"LGTG_BOT_PERMITTED_CHAT_ID":    0,
		"LGTG_BOT_POLL_INTERVAL_SECOND": 1,
		"LGTG_LINKDING_ADDRESS":         "",
		"LGTG_LINKDING_USER_TOKEN":      "",
		"LGTG_LOG_LEVEL":                "info",
		"LGTG_LOG_FILE":                 "",
	}
	if err := k.Load(confmap.Provider(defaults, "."), nil); err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}

	if err := k.Load(env.Provider("LGTG_", ".", func(s string) string {
		return s
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	ret := &Config{
		TGBotConf: struct {
			Token             string
			UpdatesBufferSize int
			PermittedChatID   int
			PollIntervalSec   int
		}{
			Token:             k.String("LGTG_BOT_TOKEN"),
			UpdatesBufferSize: k.Int("LGTG_BOT_UPDATES_BUFFER_SIZE"),
			PermittedChatID:   k.Int("LGTG_BOT_PERMITTED_CHAT_ID"),
			PollIntervalSec:   k.Int("LGTG_BOT_POLL_INTERVAL_SECOND"),
		},
		LinkdingConf: struct {
			LinkdingAddr string
			UserToken    string
		}{
			LinkdingAddr: k.String("LGTG_LINKDING_ADDRESS"),
			UserToken:    k.String("LGTG_LINKDING_USER_TOKEN"),
		},
		LogLevel: k.String("LGTG_LOG_LEVEL"),
		LogFile:  k.String("LGTG_LOG_FILE"),
	}

	if ret.LinkdingConf.LinkdingAddr == "" {
		return nil, errors.New("linkding addr is empty")
	}

	if ret.LinkdingConf.UserToken == "" {
		return nil, errors.New("linkding token is empty")
	}

	if ret.TGBotConf.Token == "" {
		return nil, errors.New("telegram bot token is empty")
	}

	return ret, nil
}
