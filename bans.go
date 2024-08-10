package etf2l

import (
	"context"
	"errors"
	"fmt"
	"github.com/leighmacdonald/steamid/v4/steamid"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"
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

	time.Sleep(time.Second * 1)

	return nextPath, nil
}

type BanOpts struct {
	Recursive
	PlayerID int    `url:"player,omitempty"` // etf2l player id only, no steamid
	Status   string `url:"status,omitempty"` // 'active' or 'expired'
	Reason   string `url:"reason,omitempty"` // 'VAC`
}

func (client *Client) Bans(ctx context.Context, httpClient HTTPExecutor, opts BanOpts) ([]Ban, error) {
	curPath := "/bans"
	max500s := 15
	cur500s := 0

	var bans []Ban

	for {
		var resp bansResponse
		if err := client.call(ctx, httpClient, curPath, opts, &resp); err != nil {
			if strings.Contains(err.Error(), "500") {
				cur500s++
				if cur500s >= max500s {
					slog.Info("Too many 500s")

					return nil, err
				}

				next, errNext := skipURLPage(curPath)
				if errNext != nil {
					return nil, errNext
				}

				curPath = next

				continue
			}

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

var errParseURL = errors.New("could not parse url")

func skipURLPage(path string) (string, error) {
	u, errParse := url.Parse(path)
	if errParse != nil {
		return "", errors.Join(errParse, errParseURL)
	}

	query := u.Query()

	page, errPage := strconv.Atoi(query.Get("page"))
	if errPage != nil {
		return "", errors.Join(errParse, errParseURL)
	}

	query.Set("page", fmt.Sprintf("%d", page+1))

	u.RawQuery = query.Encode()

	return u.String(), nil
}
