package etf2l

import (
	"context"
)

type whitelistsResponse struct {
	Status     Status               `json:"status"`
	Whitelists map[string]Whitelist `json:"whitelists"`
}

type Whitelist struct {
	Filename   string `json:"filename"`
	LastChange int    `json:"last_change"`
	URL        string `json:"url"`
}

func (client *Client) Whitelists(ctx context.Context, httpClient HTTPExecutor) (map[string]Whitelist, error) {
	var resp whitelistsResponse
	if err := client.call(ctx, httpClient, "/whitelists", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Whitelists, nil
}
