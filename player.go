package etf2l

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leighmacdonald/steamid/v4/steamid"
	"github.com/pkg/errors"
)

type PlayerResponse struct {
	Player Player `json:"player"`
	Status Status `json:"status"`
}

type BanReason struct {
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Reason string `json:"reason"`
}

type TeamCompetition struct {
	Category    string   `json:"category"`
	Competition string   `json:"competition"`
	Division    Division `json:"division"`
	URL         string   `json:"url"`
}

type IRC struct {
	Channel string `json:"channel"`
	Network string `json:"network"`
}

type SteamGroup struct {
	Avatar     string `json:"avatar"`
	SteamGroup string `json:"steam_group"`
}

type SteamPlayer struct {
	Avatar string          `json:"avatar"`
	ID     steamid.SID     `json:"id"`
	ID3    steamid.SID3    `json:"id3"`
	ID64   steamid.SteamID `json:"id64"`
}

type PlayerTeam struct {
	Competitions map[string]TeamCompetition `json:"competitions"`
	Country      string                     `json:"country"`
	Homepage     string                     `json:"homepage"`
	ID           int                        `json:"id"`
	Irc          IRC                        `json:"irc"`
	Name         string                     `json:"name"`
	Server       *string                    `json:"server"`
	Steam        SteamGroup                 `json:"steam"`
	Tag          string                     `json:"tag"`
	Type         string                     `json:"teamType"`
	Urls         URLs                       `json:"urls"`
}

type URLs struct {
	Matches   string `json:"matches"`
	Results   string `json:"results"`
	Self      string `json:"self"`
	Transfers string `json:"transfers"`
}

type PlayerClasses []string

func (f *PlayerClasses) UnmarshalJSON(data []byte) error {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch value.(type) {
	case nil, bool, []interface{}:
		*f = []string{}
	default:
		*f = value.([]string)
	}

	return nil
}

type Player struct {
	Bans       []BanReason   `json:"bans"`
	Classes    PlayerClasses `json:"classes"`
	Country    string        `json:"country"`
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Registered int           `json:"registered"`
	Steam      SteamPlayer   `json:"steam"`
	Teams      []PlayerTeam  `json:"teams"`
	Title      string        `json:"title"`
	Urls       struct {
		Results   string `json:"results"`
		Self      string `json:"self"`
		Transfers string `json:"transfers"`
	} `json:"urls"`
}

func (client *Client) Player(ctx context.Context, httpClient *http.Client, playerID string) (*Player, error) {
	var resp PlayerResponse
	if err := client.call(ctx, httpClient, fmt.Sprintf("/player/%s", playerID), nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Player, nil
}

type PlayerResultClan struct {
	Country string `json:"country"`
	Drop    bool   `json:"drop"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Steam   struct {
		Avatar string `json:"avatar"`
		Group  string `json:"group"`
	} `json:"steam"`
	URL       string `json:"url"`
	WasInTeam bool   `json:"was_in_team"`
}

type Division struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SkillContrib int    `json:"skill_contrib"`
	Tier         int    `json:"tier"`
}

type PlayerResult struct {
	Clan1       PlayerResultClan `json:"clan1"`
	Clan2       PlayerResultClan `json:"clan2"`
	Competition MatchCompetition `json:"competition"`
	Defaultwin  bool             `json:"defaultwin"`
	Division    Division         `json:"division"`
	Result      int              `json:"result"`
	Maps        []string         `json:"maps"`
	Merced      bool             `json:"merced"`
	R1          int              `json:"r1"`
	R2          int              `json:"r2"`
	Round       string           `json:"round"`
	Time        int              `json:"time"`
	Week        int              `json:"week"`
}

type links struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Active bool    `json:"active"`
}

type pagedPlayerResults struct {
	CurrentPage  int            `json:"current_page"`
	Data         []PlayerResult `json:"data"`
	FirstPageURL string         `json:"first_page_url"`
	From         int            `json:"from"`
	LastPage     int            `json:"last_page"`
	LastPageURL  string         `json:"last_page_url"`
	Links        []links        `json:"links"`
	NextPageURL  *string        `json:"next_page_url"`
	Path         string         `json:"path"`
	PerPage      int            `json:"per_page"`
	PrevPageURL  *string        `json:"prev_page_url"`
	To           int            `json:"to"`
	Total        int            `json:"total"`
}

func (resp pagedPlayerResults) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type BaseOpts struct {
	Recursive bool `url:"-"`
}

func (opts BaseOpts) IsRecursive() bool {
	return opts.Recursive
}

func (client *Client) PlayerResults(ctx context.Context, httpClient *http.Client, playerID string, opts Recursive) ([]PlayerResult, error) {
	var results []PlayerResult

	curPath := fmt.Sprintf("/player/%s/results", playerID)

	for {
		var resp pagedPlayerResults
		if err := client.call(ctx, httpClient, curPath, nil, &resp); err != nil {
			return nil, err
		}

		results = append(results, resp.Data...)

		nextURL, err := resp.NextURL(opts)
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return results, nil
}

type PlayerTransfer struct {
	By struct {
		Country string      `json:"country"`
		ID      int         `json:"id"`
		Name    string      `json:"name"`
		Steam   SteamPlayer `json:"steam"`
		URL     string      `json:"url"`
	} `json:"by"`
	Team struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Steam struct {
			Avatar string `json:"avatar"`
			Group  string `json:"group"`
		} `json:"steam"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"team"`
	Time int    `json:"time"`
	Type string `json:"type"`
}

type Meta struct {
	CurrentPage int     `json:"current_page"`
	From        int     `json:"from"`
	LastPage    int     `json:"last_page"`
	Links       []links `json:"links"`
	Path        string  `json:"path"`
	PerPage     int     `json:"per_page"`
	To          int     `json:"to"`
	Total       int     `json:"total"`
}

type transferLinks struct {
	First string  `json:"first"`
	Last  string  `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}
type playerTransfersResp struct {
	Data  []PlayerTransfer `json:"data"`
	Links transferLinks    `json:"links"`
	Meta  Meta             `json:"meta"`
}

func (resp playerTransfersResp) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Links.Next == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Links.Next)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) PlayerTransfers(ctx context.Context, httpClient *http.Client, playerID int, opts BaseOpts) ([]PlayerTransfer, error) {
	curPath := fmt.Sprintf("/player/%d/transfers", playerID)

	var transfers []PlayerTransfer

	for {
		var resp playerTransfersResp
		if err := client.call(ctx, httpClient, curPath, nil, &resp); err != nil {
			return nil, err
		}

		transfers = append(transfers, resp.Data...)

		nextURL, err := resp.NextURL(opts)
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return transfers, nil
}
