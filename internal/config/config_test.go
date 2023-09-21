package config

import (
	"testing"

	"github.com/c2fo/testify/require"
)

func TestConfig(t *testing.T) {
	lnkgAddr := "192.168.1.1"
	botToken := "AFBHASBFAJKFNAKBFAHFA"

	t.Setenv("LGTG_BOT_TOKEN", botToken)
	t.Setenv("LGTG_LINKDING_ADDRESS", lnkgAddr)

	config, err := GetConfig()
	require.NoError(t, err)

	require.Equal(t, config.LinkdingConf.LinkdingAddr, lnkgAddr)
	require.Equal(t, config.TGBitConf.Token, botToken)
}
