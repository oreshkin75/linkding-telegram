package telegram

import (
	"linkding-telegram/internal/utils"
	"unicode/utf16"
)

func (m Message) ExtractURLs() []string {
	var urls []string
	seen := make(map[string]struct{})

	addURL := func(rawURL string) {
		if rawURL == "" {
			return
		}
		if _, ok := seen[rawURL]; ok {
			return
		}

		seen[rawURL] = struct{}{}
		urls = append(urls, rawURL)
	}

	addURL(m.LinkPrev.URL)
	addEntityURLs(m.Text, m.Entities, addURL)
	addEntityURLs(m.Caption, m.CaptionEntities, addURL)

	for _, rawURL := range utils.ParseURLs(m.Text) {
		addURL(rawURL)
	}
	for _, rawURL := range utils.ParseURLs(m.Caption) {
		addURL(rawURL)
	}

	return urls
}

func addEntityURLs(text string, entities []MessageEntity, addURL func(string)) {
	for _, entity := range entities {
		switch entity.Type {
		case "text_link":
			addURL(entity.URL)
		case "url":
			addURL(substringByUTF16Offset(text, entity.Offset, entity.Length))
		}
	}
}

func substringByUTF16Offset(text string, offset, length int64) string {
	if offset < 0 || length <= 0 {
		return ""
	}

	encoded := utf16.Encode([]rune(text))
	start := int(offset)
	end := start + int(length)
	if start < 0 || start >= len(encoded) || end > len(encoded) {
		return ""
	}

	return string(utf16.Decode(encoded[start:end]))
}
