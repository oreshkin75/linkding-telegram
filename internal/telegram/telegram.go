package telegram

import (
	"fmt"
	"net/url"
	"strings"

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

func (t *Telegram) PollingUpdates() error {
	api := echotron.NewAPI(t.tgToken)

	for u := range echotron.PollingUpdates(t.tgToken) {
		if u.Message != nil {
			if u.Message.Text != "" {
				words := strings.Fields(u.Message.Text)
				var links []string

				for _, word := range words {
					u, err := url.Parse(word)
					if err != nil {
						continue
					}
					if u.Scheme == "http" || u.Scheme == "https" {
						links = append(links, word)
					}
				}

				fmt.Println("==", links)
			}
			if u.Message.CaptionEntities != nil {
				fmt.Println("---", len(u.Message.CaptionEntities))
			}
			if u.Message.Entities != nil {
				for _, msg := range u.Message.Entities {
					if msg.URL != "" {
						fmt.Println(msg.URL)
					}
				}
			}
		}
		if u.Message.Text == "/start" {
			api.SendMessage("Hello world", u.ChatID(), nil)
		}
	}

	return nil
}
