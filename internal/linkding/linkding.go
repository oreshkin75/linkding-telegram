package linkding

import (
	"context"
	"errors"
	"io"
	"linkding-telegram/internal/config"
	"net/http"
	"net/url"
)

type Linkding struct {
	addr    string
	lgToken string
}

func NewLinkding(opts config.LinkdingConf) *Linkding {
	return &Linkding{
		addr:    opts.LinkdingAddr,
		lgToken: opts.UserToken,
	}
}

func (l *Linkding) GetBookmarks(ctx context.Context, filters, limit, offset string) ([]byte, error) {
	baseURL, err := url.Parse(l.addr + "/api/bookmarks")
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse url"), err)
	}

	params := url.Values{}
	if filters != "" {
		params.Add("q", filters)
	}
	if limit != "" {
		params.Add("limit", limit)
	}
	if offset != "" {
		params.Add("offset", offset)
	}

	baseURL.RawQuery = params.Encode()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create request"), err)
	}

	req.Header.Add("Authorization", "Token "+l.lgToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("failed to send request"), err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to read body"), err)
	}

	return body, nil
}
