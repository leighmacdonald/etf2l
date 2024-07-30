package etf2l

import (
	"context"
	"net/http"

	"github.com/leighmacdonald/steamid/v4/steamid"
	"github.com/pkg/errors"
)

type Ban struct {
	Start     int             `json:"start"`
	End       int             `json:"end"`
	Name      string          `json:"name"`
	Steamid   string          `json:"steamid"`
	Steamid64 steamid.SteamID `json:"steamid64"`
	Profile   string          `json:"profile"`
	Expired   bool            `json:"expired"`
	Reason    string          `json:"reason"`
}

type pagedBans struct {
	CurrentPage  int     `json:"current_page"`
	Data         []Ban   `json:"data"`
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

type bansResponse struct {
	Pager pagedBans `json:"bans"`
}

func (resp bansResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type BanOpts struct {
	Recursive
	PlayerID int    `url:"player,omitempty"` // etf2l player id only, no steamid
	Status   string `url:"status,omitempty"` // 'active' or 'expired'
	Reason   string `url:"reason,omitempty"` // 'VAC`
}

func (client *Client) Bans(ctx context.Context, httpClient *http.Client, opts BanOpts) ([]Ban, error) {
	curPath := "/bans"

	var bans []Ban

	for {
		// TODO Remove, for some reason this page 500s.
		if curPath == "https://api-v2.etf2l.org/bans?page=18" {
			curPath = "https://api-v2.etf2l.org/bans?page=19"
		}

		var resp bansResponse
		if err := client.call(ctx, httpClient, curPath, opts, &resp); err != nil {
			return nil, err
		}

		bans = append(bans, resp.Pager.Data...)

		nextURL, err := resp.NextURL(opts)

		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return bans, nil
}
