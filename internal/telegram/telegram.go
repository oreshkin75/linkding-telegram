package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"linkding-telegram/internal/config"
	"linkding-telegram/internal/linkding"
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
	tgApiUrlWithToken string
	permittedChatIDs  []int
	pollIntervalSec   int
	linkding          *linkding.Linkding
	log               *logrus.Logger
	updates           chan Update
}

func New(opts *config.Config, linkding *linkding.Linkding, log *logrus.Logger) *Telegram {
	return &Telegram{
		tgApiUrlWithToken: tgApiURL + opts.TGBotConf.Token,
		log:               log,
		linkding:          linkding,
		permittedChatIDs:  opts.TGBotConf.PermittedChatIDs,
		pollIntervalSec:   opts.TGBotConf.PollIntervalSec,
		updates:           make(chan Update, opts.TGBotConf.UpdatesBufferSize),
	}
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

func (t *Telegram) getUpdates(ctx context.Context, offset string) (Response, error) {
	baseURL, err := url.Parse(t.tgApiUrlWithToken + getUpdatesMethod)
	if err != nil {
		return Response{}, fmt.Errorf("failed to parse get updates url: %w", err)
	}

	params := url.Values{}
	params.Add("offset", offset)

	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return Response{}, fmt.Errorf("failed to make get updates GET request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("failed to create http client: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read get updated body: %w", err)
	}

	var restResponse Response

	if err := json.Unmarshal(body, &restResponse); err != nil {
		return Response{}, fmt.Errorf("failed to unmarshal get updates body: %w", err)
	}

	return restResponse, nil
}

func (t *Telegram) PollUpdates(ctx context.Context) {
	var lastUpdateID int
	for {
		select {
		case <-ctx.Done():
			t.log.Info("receiving updates has been stopped")
			return
		default:
			response, err := t.getUpdates(ctx, fmt.Sprintf("%d", lastUpdateID+1))
			if err != nil {
				t.log.Error(err)
				continue
			}

			for _, update := range response.Result {
				if !t.checkChatID(int(update.Message.Chat.ID)) {
					t.log.WithFields(logrus.Fields{
						"chatID": update.UpdateID,
					}).Info("message received from an unauthorized user")

					lastUpdateID = update.UpdateID

					continue
				}

				t.updates <- update
				lastUpdateID = update.UpdateID
			}

			time.Sleep(time.Second * time.Duration(t.pollIntervalSec))
		}
	}
}

func (t *Telegram) GetUpdate() <-chan Update {
	return t.updates
}

func (t *Telegram) checkChatID(chatID int) bool {
	if len(t.permittedChatIDs) == 0 {
		return true
	}

	for _, permittedChatID := range t.permittedChatIDs {
		if permittedChatID == chatID {
			return true
		}
	}

	return false
}
