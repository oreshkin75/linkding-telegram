package main

import (
	"fmt"
	"linkding-telegram/internal/config"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetLevel(log.WarnLevel)
}

func setupLogger(logLevel string, isLogToFile bool) error {
	const logPath = "log.txt"

	if isLogToFile {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open/create log file: %w", err)
		}

		log.SetOutput(file)
	}

	switch logLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	}

	return nil
}

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to read configuration file: %w", err)
	}

	fmt.Println(config)

}
