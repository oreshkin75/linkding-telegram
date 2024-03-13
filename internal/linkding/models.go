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

type CreateBookmarkResp struct {
	ID                 int      `json:"id"`
	URL                string   `json:"url"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Notes              string   `json:"notes"`
	WebsiteTitle       string   `json:"website_title"`
	WebsiteDescription string   `json:"website_description"`
	IsArchived         bool     `json:"is_archived"`
	Unread             bool     `json:"unread"`
	Shared             bool     `json:"shared"`
	TagNames           []string `json:"tag_names"`
	DateAdded          string   `json:"date_added"`
	DateModified       string   `json:"date_modified"`
}
