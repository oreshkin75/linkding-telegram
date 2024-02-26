package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"linkding-telegram/internal/config"
	"linkding-telegram/internal/linkding"
	"linkding-telegram/internal/utils"
	"net/http"
	"net/url"
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
	linkding          *linkding.Linkding
	log               *logrus.Logger
}

func New(opts config.TGBotConf, linkding *linkding.Linkding, log *logrus.Logger) *Telegram {
	return &Telegram{
		botToken:          opts.Token,
		tgApiUrlWithToken: tgApiURL + opts.Token,
		log:               log,
		linkding:          linkding,
	}
}

func (t *Telegram) getUpdates(ctx context.Context, offset string) ([]Update, error) {
	baseURL, err := url.Parse(t.tgApiUrlWithToken + getUpdatesMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to parse get updates url: %w", err)
	}

	params := url.Values{}
	params.Add("offset", offset)

	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
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
			updates, err := t.getUpdates(ctx, fmt.Sprintf("%d", lastUpdateID+1))
			if err != nil {
				t.log.Error("failed to get updates from telegram")
				continue
			}

			for _, update := range updates {
				if update.Message.Text != "" {
					links := utils.ParseURLs(update.Message.Text)
					t.log.WithFields(logrus.Fields{ // TODO test
						"links": links,
					}).Info("links in message")

					for _, link := range links {
						resp, err := t.linkding.CreateBookmark(ctx, &linkding.CreateBookmark{
							URL:    link,
							Unread: true,
						})

						if err != nil {
							t.log.Error(err)
							continue
						}

						fmt.Println(string(resp))
					}
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
