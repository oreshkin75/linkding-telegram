package config

import (
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestConfig(t *testing.T) {
	testCase := []struct {
		name  string
		input struct {
			linkdingAddr string
			botToken     string
			logLevel     string
		}
		expected struct {
			linkdingAddr string
			botToken     string
			logLevel     string
		}
	}{{
		name: "Simple test",
		input: struct {
			linkdingAddr string
			botToken     string
			logLevel     string
		}{
			linkdingAddr: "192.168.1.1",
			botToken:     "AFBHASBFAJKFNAKBFAHFA",
			logLevel:     "info",
		},
		expected: struct {
			linkdingAddr string
			botToken     string
			logLevel     string
		}{
			linkdingAddr: "192.168.1.1",
			botToken:     "AFBHASBFAJKFNAKBFAHFA",
			logLevel:     "info",
		},
	},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("LGTG_BOT_TOKEN", tc.input.botToken)
			t.Setenv("LGTG_LINKDING_ADDRESS", tc.input.linkdingAddr)
			t.Setenv("LGTG_LOG_LEVEL", tc.input.logLevel)

			config, err := GetConfig()
			assert.NoError(t, err)

			assert.Equal(t, config.LinkdingConf.LinkdingAddr, tc.expected.linkdingAddr)
			assert.Equal(t, config.TGBitConf.Token, tc.expected.botToken)
			assert.Equal(t, config.LogLevel, tc.expected.logLevel)
		})
	}
}
