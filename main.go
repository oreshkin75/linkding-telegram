package main

import (
	"context"
	"fmt"
	"linkding-telegram/internal/config"
	"linkding-telegram/internal/linkding"
	"linkding-telegram/internal/telegram"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetOutput(os.Stdout)

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetLevel(logrus.InfoLevel)
}

func setupLogger(logLevel, logPath string) error {
	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666) //nolint:gosec // log path is intentionally configured by the service operator.
		if err != nil {
			return fmt.Errorf("failed to open/create log file: %w", err)
		}

		log.SetOutput(file)
	}

	switch logLevel {
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	}

	return nil
}

func main() {
	appConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to read configuration file: %s", err.Error())
	}

	if err := setupLogger(appConfig.LogLevel, appConfig.LogFile); err != nil {
		log.Fatalf("failed to setup logger: %s", err.Error())
	}

	lnkdng := linkding.New(appConfig, log)
	tg := telegram.New(appConfig, lnkdng, log)

	ctx := context.Background()

	go tg.PollUpdates(ctx)

	for update := range tg.GetUpdate() {
		processUpdate(ctx, update, appConfig, lnkdng)
	}
}

func processUpdate(ctx context.Context, update telegram.Update, appConfig *config.Config, lnkdng *linkding.Linkding) {
	urls := update.Message.ExtractURLs()
	if len(urls) == 0 {
		log.WithFields(logrus.Fields{
			"messageID": update.Message.MessageID,
			"chatID":    update.Message.Chat.ID,
		}).Warn("message does not contain URLs")

		return
	}

	tags := bookmarkTags(update.Message, appConfig.LinkdingConf.DefaultTag)
	for _, rawURL := range urls {
		req := bookmarkRequest(rawURL, tags)
		if _, err := lnkdng.CreateBookmark(ctx, req); err != nil {
			log.Fatalf("failed to create bookmark in Linkding: %s", err.Error())
		}
	}
}

func bookmarkTags(message telegram.Message, defaultTag string) []string {
	var tags []string
	if message.MessageOrigin.Chat.Username != "" {
		tags = append(tags, message.MessageOrigin.Chat.Username)
	}

	if defaultTag != "" {
		tags = append(tags, defaultTag)
	}

	return tags
}

func bookmarkRequest(rawURL string, tags []string) *linkding.CreateBookmarkReqBody {
	req := &linkding.CreateBookmarkReqBody{
		URL:    rawURL,
		Unread: true,
	}
	if len(tags) > 0 {
		req.TagNames = tags
	}

	return req
}
