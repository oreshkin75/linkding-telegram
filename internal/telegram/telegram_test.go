package telegram

import "testing"

func TestTelegram(t *testing.T) {
	// TODO delete
	tgToken := "6386280311:AAEwAngQVokrfm2IurwjEUHTRV2BcLBCDeI"

	tg := NewTelegram(tgToken)
	tg.PollingUpdates()
}