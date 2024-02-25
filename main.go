package main

import (
	"context"
	"fmt"
	"linkding-telegram/internal/config"
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
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to read configuration file: %s", err.Error())
	}

	if err := setupLogger(config.LogLevel, config.LogFile); err != nil {
		log.Fatalf("failed to setup logger: %s", err.Error())
	}

	tg := telegram.New(config.TGBitConf.Token, log)

	tg.PollUpdates(context.Background())
}
