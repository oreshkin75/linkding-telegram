package config

import (
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		name          string
		linkdingAddr  string
		linkdingToken string
		botToken      string
		logLevel      string
	}{
		{
			name:          "reads current environment variables",
			linkdingAddr:  "http://192.168.1.1:9090",
			linkdingToken: "linkding-token",
			botToken:      "telegram-token",
			logLevel:      "info",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("BOT_TOKEN", tc.botToken)
			t.Setenv("LINKDING_ADDRESS", tc.linkdingAddr)
			t.Setenv("LINKDING_USER_TOKEN", tc.linkdingToken)
			t.Setenv("LOG_LEVEL", tc.logLevel)

			config, err := GetConfig()
			assert.NoError(t, err)

			assert.Equal(t, tc.linkdingAddr, config.LinkdingConf.LinkdingAddr)
			assert.Equal(t, tc.linkdingToken, config.LinkdingConf.UserToken)
			assert.Equal(t, tc.botToken, config.TGBotConf.Token)
			assert.Equal(t, tc.logLevel, config.LogLevel)
		})
	}
}
