package telegram

import (
	"encoding/json"
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestMessageExtractURLs(t *testing.T) {
	testCases := []struct {
		name     string
		message  Message
		expected []string
	}{
		{
			name: "uses link preview URL",
			message: Message{
				LinkPrev: LinkPrev{URL: "https://example.com/preview"},
			},
			expected: []string{"https://example.com/preview"},
		},
		{
			name: "falls back to text URL when preview is absent",
			message: Message{
				Text: "save https://example.com/from-text please",
			},
			expected: []string{"https://example.com/from-text"},
		},
		{
			name: "extracts text_link entity URL",
			message: Message{
				Text: "click here",
				Entities: []MessageEntity{
					{Type: "text_link", Offset: 0, Length: 10, URL: "https://example.com/text-link"},
				},
			},
			expected: []string{"https://example.com/text-link"},
		},
		{
			name: "extracts URL entity using UTF-16 offsets",
			message: Message{
				Text: "🔥 https://example.com/entity",
				Entities: []MessageEntity{
					{Type: "url", Offset: 3, Length: 26},
				},
			},
			expected: []string{"https://example.com/entity"},
		},
		{
			name: "extracts caption URLs",
			message: Message{
				Caption: "caption https://example.com/caption",
			},
			expected: []string{"https://example.com/caption"},
		},
		{
			name: "deduplicates URLs keeping first occurrence order",
			message: Message{
				LinkPrev: LinkPrev{URL: "https://example.com/same"},
				Text:     "https://example.com/same https://example.com/other",
			},
			expected: []string{"https://example.com/same", "https://example.com/other"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.message.ExtractURLs())
		})
	}
}

func TestMessageUnmarshalEntitiesAndCaptionEntities(t *testing.T) {
	payload := []byte(`{
		"message_id": 1,
		"text": "https://example.com/text",
		"entities": [{"type": "url", "offset": 0, "length": 24}],
		"caption": "caption https://example.com/caption",
		"caption_entities": [{"type": "url", "offset": 8, "length": 27}]
	}`)

	var message Message
	err := json.Unmarshal(payload, &message)

	assert.NoError(t, err)
	assert.Len(t, message.Entities, 1)
	assert.Len(t, message.CaptionEntities, 1)
}
