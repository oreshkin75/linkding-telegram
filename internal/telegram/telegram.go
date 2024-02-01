package telegram

import (
	"fmt"
	"linkding-telegram/internal/utils"

	"github.com/NicoNex/echotron/v3"
)

type Telegram struct {
	tgToken string
}

func NewTelegram(tgToken string) *Telegram {
	return &Telegram{
		tgToken: tgToken,
	}
}

func (t *Telegram) PollingUrlUpdates() {
	for u := range echotron.PollingUpdates(t.tgToken) {
		if u.Message != nil {
			if u.Message.Text != "" {
				links := utils.ParseURLs(u.Message.Text)

				fmt.Println("==", links)
			}
			if u.Message.CaptionEntities != nil {
				for _, msg := range u.Message.CaptionEntities {
					if msg.URL != "" {
						fmt.Println(msg.URL)
					}
				}
			}
			if u.Message.Entities != nil {
				for _, msg := range u.Message.Entities {
					if msg.URL != "" {
						fmt.Println(msg.URL)
					}
				}
			}
		}
	}
}

func (t *Telegram) SelectLink(chatID int64, s []string) {
	api := echotron.NewAPI(t.tgToken)
	api.SendMessage("Select the desired link", chatID, nil)
}
