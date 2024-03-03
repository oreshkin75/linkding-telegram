package linkding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"linkding-telegram/internal/config"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

const (
	bookmarksMethods = "/api/bookmarks/"
)

type Linkding struct {
	addr    string
	lgToken string
	log     *logrus.Logger
}

func New(opts *config.Config, log *logrus.Logger) *Linkding {
	return &Linkding{
		addr:    opts.LinkdingConf.LinkdingAddr,
		lgToken: opts.LinkdingConf.UserToken,
		log:     log,
	}
}

func (l *Linkding) GetBookmarks(ctx context.Context, filters, limit, offset string) ([]byte, error) {
	baseURL, err := url.Parse(l.addr + bookmarksMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
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
		return nil, fmt.Errorf("failed to create get bookmarks GET request: %w", err)
	}

	req.Header.Add("Authorization", "Token "+l.lgToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

func (l *Linkding) CreateBookmark(ctx context.Context, opts *CreateBookmark) ([]byte, error) {
	reqBytes, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create bookmark body: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", l.addr+bookmarksMethods, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create create bookmark POST request: %w", err)
	}

	req.Header.Add("Authorization", "Token "+l.lgToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		l.log.Errorf("status code of POST %s - %d", l.addr+bookmarksMethods, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}
