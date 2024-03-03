package telegram

type Response struct {
	OK          bool     `json:"ok"`
	Description string   `json:"description,omitempty"`
	Result      []Update `json:"result,omitempty"`
	ErrCode     int64    `json:"error_code,omitempty"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message,omitempty"`
}

type Message struct {
	MessageID       int64                `json:"message_id"`
	SenderChat      Chat                 `json:"sender_chat,omitempty"`
	Chat            Chat                 `json:"chat,omitempty"`
	Text            string               `json:"text,omitempty"`
	MessageEntities []MessageEntity      `json:"caption_entities,omitempty"`
	MessageOrigin   MessageOriginChannel `json:"forward_origin,omitempty"`
	LinkPrev        LinkPrev             `json:"link_preview_options,omitempty"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	Url    string `json:"url,omitempty"`
}

type MessageOriginChannel struct {
	Type      string `json:"type,omitempty"`
	Date      int64  `json:"date,omitempty"`
	Chat      Chat   `json:"chat,omitempty"`
	MessageID int64  `json:"message_id,omitempty"`
}

type LinkPrev struct {
	Url string `json:"url,omitempty"`
}

type SendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
