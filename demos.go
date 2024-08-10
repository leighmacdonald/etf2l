package etf2l

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

type Demo struct {
	ID          int    `json:"id"`
	Time        int    `json:"time"`
	Match       int    `json:"match"`
	DownloadURL string `json:"download_url"`
	Stv         bool   `json:"stv"`
	FirstPerson bool   `json:"first_person"`
	Downloads   int    `json:"downloads"`
	Owner       int    `json:"owner"`
	OwnerName   string `json:"owner_name"`
	Pruned      bool   `json:"pruned"`
	File        string `json:"file"`
	Extension   string `json:"extension"`
}

type DemoPage struct {
	CurrentPage  int     `json:"current_page"`
	Data         []Demo  `json:"data"`
	FirstPageURL string  `json:"first_page_url"`
	From         int     `json:"from"`
	LastPage     int     `json:"last_page"`
	LastPageURL  string  `json:"last_page_url"`
	Links        []links `json:"links"`
	NextPageURL  *string `json:"next_page_url"`
	Path         string  `json:"path"`
	PerPage      int     `json:"per_page"`
	PrevPageURL  *string `json:"prev_page_url"`
	To           int     `json:"to"`
	Total        int     `json:"total"`
}

type demosResponse struct {
	Status Status   `json:"status"`
	Pager  DemoPage `json:"demos"`
}

func (resp demosResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type DemoOpts struct {
	Recursive
	PlayerID string   `url:"player,omitempty"`
	Type     []string `url:"type,omitempty"` // stv, first_person
	Pruned   bool     `url:"pruned,omitempty"`
	From     int      `url:"from,omitempty"` // unixtime start
	To       int      `url:"to,omitempty"`   // unixtime end
}

func getPath(path string) (string, error) {
	parsed, err := url.ParseRequestURI(path)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse URL")
	}

	return parsed.Path + "?" + parsed.RawQuery, nil
}

func (client *Client) Demos(ctx context.Context, httpClient HTTPExecutor, opts Recursive) ([]Demo, error) {
	var demos []Demo

	curPath := "/demos"

	for {
		var resp demosResponse
		if err := client.call(ctx, httpClient, curPath, opts, &resp); err != nil {
			return nil, err
		}

		demos = append(demos, resp.Pager.Data...)

		nextURL, err := resp.NextURL(opts)
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return demos, nil
}
