package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"linkding-telegram/internal/utils"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	tgApiURL          = "https://api.telegram.org/bot"
	getUpdatesMethod  = "/getUpdates"
	sendMessageMethod = "/sendMessage"
)

type Telegram struct {
	botToken          string
	tgApiUrlWithToken string
	log               *logrus.Logger
}

func New(botToken string, log *logrus.Logger) *Telegram {
	return &Telegram{
		botToken:          botToken,
		tgApiUrlWithToken: tgApiURL + botToken,
		log:               log,
	}
}

func (t *Telegram) getUpdates(ctx context.Context, offset int) ([]Update, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s?offset=%d", t.tgApiUrlWithToken, getUpdatesMethod, offset), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make get updates GET request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read get updated body: %w", err)
	}

	var restResponse struct {
		Result []Update `json:"result"`
	}

	if err := json.Unmarshal(body, &restResponse); err != nil {
		return nil, fmt.Errorf("faile to unmarshal get updates body: %w", err)
	}

	return restResponse.Result, nil
}

func (t *Telegram) sendMessage(ctx context.Context, chatID int64, text string) error {
	reqBody := &SendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal send message body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.tgApiUrlWithToken+sendMessageMethod, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("failed to make send message POST request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	if _, err = http.DefaultClient.Do(req); err != nil {
		return fmt.Errorf("failed to create http client: %w", err)
	}

	return nil
}

func (t *Telegram) PollUpdates(ctx context.Context) {
	var lastUpdateID int
	for {
		select {
		case <-ctx.Done():
			t.log.Info("receiving updates has been stopped")
			return
		default:
			updates, err := t.getUpdates(ctx, lastUpdateID+1)
			if err != nil {
				t.log.Error("failed to get updates from telegram")
				continue
			}

			for _, update := range updates {
				t.log.WithFields(logrus.Fields{ // TODO test
					"message": update.Message.Text,
					"chat_id": update.Message.Chat.ID,
				}).Info("tg message")

				if update.Message.Text != "" {
					links := utils.ParseURLs(update.Message.Text)
					t.log.WithFields(logrus.Fields{ // TODO test
						"links": links,
					}).Info("links in message")
				}

				if err := t.sendMessage(ctx, update.Message.Chat.ID, "Получено: "+update.Message.Text); err != nil { // TODO test
					t.log.Error("failed to send message to telegram")
					continue
				}
				lastUpdateID = update.UpdateID
			}

			time.Sleep(1 * time.Second)
		}
	}
}

/*type Telegram struct {
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
}*/
