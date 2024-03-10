package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	TGBotConf    TGBotConf    `envPrefix:"BOT_"`
	LinkdingConf LinkdingConf `envPrefix:"LINKDING_"`
	LogLevel     string       `env:"LOG_LEVEL" envDefault:"info"`
	LogFile      string       `env:"LOG_FILE"`
}

type TGBotConf struct {
	Token             string `env:"TOKEN,notEmpty"`
	UpdatesBufferSize int    `env:"UPDATES_BUFFER_SIZE" envDefault:"1"`
	PermittedChatIDs  []int  `env:"PERMITTED_CHAT_IDS" envSeparator:","`
	PollIntervalSec   int    `env:"POLL_INTERVAL_SECOND" envDefault:"1"`
}

type LinkdingConf struct {
	LinkdingAddr string `env:"ADDRESS,notEmpty"`
	UserToken    string `env:"USER_TOKEN,notEmpty"`
	DefaultTag   string `env:"DEFAULT_TAG"`
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
