package config

import (
	"errors"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

type Config struct {
	TGBitConf    TGBotConf
	LinkdingConf LinkdingConf
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
		return nil, errors.Join(errors.New("failed to load configuration"), err)
	}

	return &Config{
		TGBitConf: TGBotConf{
			Token: k.String("LGTG_BOT_TOKEN"),
		},
		LinkdingConf: LinkdingConf{
			LinkdingAddr: k.String("LGTG_LINKDING_ADDRESS"),
			UserToken:    k.String("LGTG_LINKDING_USER_TOKEN"),
		},
	}, nil
}
