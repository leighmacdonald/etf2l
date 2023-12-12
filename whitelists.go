package etf2l

import (
	"context"
	"net/http"
)

type WhitelistsResponse struct {
	Status     Status               `json:"status"`
	Whitelists map[string]Whitelist `json:"whitelists"`
}

type Whitelist struct {
	Filename   string `json:"filename"`
	LastChange int    `json:"last_change"`
	Url        string `json:"url"`
}

func (client *Client) Whitelists(ctx context.Context) (map[string]Whitelist, error) {
	var resp WhitelistsResponse
	if err := client.call(ctx, http.MethodGet, "/whitelists", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Whitelists, nil
}
