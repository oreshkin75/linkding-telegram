package linkding

type CreateBookmarkReqBody struct {
	URL         string   `json:"url"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	IsArchived  bool     `json:"is_archived,omitempty"`
	Unread      bool     `json:"unread,omitempty"`
	Shared      bool     `json:"shared,omitempty"`
	TagNames    []string `json:"tag_names,omitempty"`
}
