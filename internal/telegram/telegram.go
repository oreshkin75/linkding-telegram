package telegram

import (
	"fmt"
	"linkding-telegram/internal/utils"

	"github.com/NicoNex/echotron/v3"
)

type Telegram struct {
	tgToken string
}

func New(tgToken string) *Telegram {
	return &Telegram{
		tgToken: tgToken,
	}
}

func (t *Telegram) PollingURLUpdates() {
	for u := range echotron.PollingUpdates(t.tgToken) {
		if u.Message == nil {
			return
		}

		if u.Message.Text != "" {
			links := utils.ParseURLs(u.Message.Text)

			fmt.Println("==", links) // TODO for test purposes
		}

		if u.Message.CaptionEntities != nil {
			for _, msg := range u.Message.CaptionEntities {
				if msg.URL != "" {
					fmt.Println(msg.URL) // TODO test purposes only
				}
			}
		}

		if u.Message.Entities != nil {
			for _, msg := range u.Message.Entities {
				if msg.URL != "" {
					fmt.Println(msg.URL) // TODO for test purpose
				}
			}
		}

	}
}

func (t *Telegram) SelectLink(chatID int64, s []string) error {
	api := echotron.NewAPI(t.tgToken)
	_, err := api.SendMessage("Select the desired link", chatID, nil)
	if err != nil {
		return fmt.Errorf("failed to send tg message: %w", err)
	}

	return nil
}
