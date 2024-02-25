package config

import (
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

type Config struct {
	TGBitConf    TGBotConf
	LinkdingConf LinkdingConf
	LogLevel     string
	LogFile      string
}

type TGBotConf struct {
	Token string
}

type LinkdingConf struct {
	LinkdingAddr string
	UserToken    string
}

func GetConfig() (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(env.Provider("LGTG_", ".", func(s string) string {
		return s
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &Config{
		TGBitConf: TGBotConf{
			Token: k.String("LGTG_BOT_TOKEN"),
		},
		LinkdingConf: LinkdingConf{
			LinkdingAddr: k.String("LGTG_LINKDING_ADDRESS"),
			UserToken:    k.String("LGTG_LINKDING_USER_TOKEN"),
		},
		LogLevel: k.String("LGTG_LOG_LEVEL"),
		LogFile:  k.String("LGTG_LOG_FILE"),
	}, nil
}
