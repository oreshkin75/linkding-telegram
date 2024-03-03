package utils

import (
	"regexp"
)

func ParseURLs(str string) []string {
	urlRegex := regexp.MustCompile(`\bhttps?://[^\s]+\b`)
	return urlRegex.FindAllString(str, -1)
}
