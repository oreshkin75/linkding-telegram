package utils

import (
	"net/url"
	"strings"
)

func ParseURLs(str string) []string {
	words := strings.Fields(str)
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

	return links
}
