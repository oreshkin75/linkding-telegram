package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	TGBotConf    TGBotConf
	LinkdingConf LinkdingConf
	LogLevel     string
	LogFile      string
}

type TGBotConf struct {
	Token             string
	UpdatesBufferSize int
	PermittedChatIDs  []int
	PollIntervalSec   int
}

type LinkdingConf struct {
	LinkdingAddr string
	UserToken    string
}

func GetConfig() (*Config, error) {
	viper.SetEnvPrefix("LGTG")
	viper.AutomaticEnv()

	viper.SetDefault("BOT_TOKEN", "")
	viper.SetDefault("BOT_UPDATES_BUFFER_SIZE", 1)
	viper.SetDefault("BOT_PERMITTED_CHAT_IDS", []int{})
	viper.SetDefault("BOT_POLL_INTERVAL_SECOND", 1)
	viper.SetDefault("LINKDING_ADDRESS", "")
	viper.SetDefault("LINKDING_USER_TOKEN", "")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE", "")

	fmt.Println("---", viper.GetStringSlice("BOT_PERMITTED_CHAT_ID"))
	config := &Config{
		TGBotConf: TGBotConf{
			Token:             viper.GetString("BOT_TOKEN"),
			UpdatesBufferSize: viper.GetInt("BOT_UPDATES_BUFFER_SIZE"),
			PermittedChatIDs:  viper.GetIntSlice("BOT_PERMITTED_CHAT_ID"),
			PollIntervalSec:   viper.GetInt("BOT_POLL_INTERVAL_SECOND"),
		},
		LinkdingConf: LinkdingConf{
			LinkdingAddr: viper.GetString("LINKDING_ADDRESS"),
			UserToken:    viper.GetString("LINKDING_USER_TOKEN"),
		},
		LogLevel: viper.GetString("LOG_LEVEL"),
		LogFile:  viper.GetString("LOG_FILE"),
	}

	fmt.Println("====", config.TGBotConf.PermittedChatIDs, len(config.TGBotConf.PermittedChatIDs))
	fmt.Println(config)
	return config, nil
}
