package etf2l

import (
	"context"
	"net/http"
)

type Demo struct {
	Id          int    `json:"id"`
	Time        int    `json:"time"`
	Match       int    `json:"match"`
	DownloadUrl string `json:"download_url"`
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
	CurrentPage  int    `json:"current_page"`
	Data         []Demo `json:"data"`
	FirstPageUrl string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageUrl  string `json:"last_page_url"`
	Links        []struct {
		Url    *string `json:"url"`
		Label  string  `json:"label"`
		Active bool    `json:"active"`
	} `json:"links"`
	NextPageUrl string      `json:"next_page_url"`
	Path        string      `json:"path"`
	PerPage     int         `json:"per_page"`
	PrevPageUrl interface{} `json:"prev_page_url"`
	To          int         `json:"to"`
	Total       int         `json:"total"`
}

type DemosResponse struct {
	Status Status   `json:"status"`
	Pager  DemoPage `json:"demos"`
}

type DemoOpts struct {
	PlayerID string   `json:"player,omitempty"`
	Type     []string `json:"type,omitempty"` // stv, first_person
	Pruned   bool     `json:"pruned,omitempty"`
	From     int      `json:"from,omitempty"` // unixtime start
	To       int      `json:"to,omitempty"`   // unixtime end
}

func (client *Client) Demos(ctx context.Context, opts DemoOpts) ([]Demo, error) {
	var resp DemosResponse
	if err := client.call(ctx, http.MethodGet, "/demos", opts, &resp); err != nil {
		return nil, err
	}

	return resp.Pager.Data, nil
}
