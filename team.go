package etf2l

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Team struct {
	Competitions map[string]TeamCompetition `json:"competitions"`
	Country      string                     `json:"country"`
	Homepage     string                     `json:"homepage"`
	ID           int                        `json:"id"`
	Irc          struct {
		Channel string `json:"channel"`
		Network string `json:"network"`
	} `json:"irc"`
	Name   string     `json:"name"`
	Server string     `json:"server"`
	Steam  SteamGroup `json:"steam"`
	Tag    string     `json:"tag"`
	Urls   struct {
		Matches   string `json:"matches"`
		Results   string `json:"results"`
		Self      string `json:"self"`
		Transfers string `json:"transfers"`
	} `json:"urls"`
	Players     []TeamPlayer     `json:"players"`
	NameChanges []TeamNameChange `json:"name_changes"`
}
type TeamNameChange struct {
	From string `json:"from"`
	To   string `json:"to"`
	Time int    `json:"time"`
}
type TeamPlayer struct {
	Country string      `json:"country"`
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Role    string      `json:"role"`
	Steam   SteamPlayer `json:"steam"`
	URL     string      `json:"url"`
}
type teamResponse struct {
	Team   Team   `json:"team"`
	Status Status `json:"status"`
}

func (client *Client) Team(ctx context.Context, teamID int) (*Team, error) {
	var resp teamResponse
	if err := client.call(ctx, fmt.Sprintf("/team/%d", teamID), nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Team, nil
}

type TransferPlayerInfo struct {
	Country string      `json:"country"`
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Steam   SteamPlayer `json:"steam"`
	URL     string      `json:"url"`
}

type TeamTransfer struct {
	Who  TransferPlayerInfo `json:"who"`
	By   TransferPlayerInfo `json:"by"`
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

type teamTransfersResp struct {
	Data  []TeamTransfer `json:"data"`
	Meta  Meta           `json:"meta"`
	Links transferLinks  `json:"links"`
}

func (resp teamTransfersResp) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Links.Next == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Links.Next)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) TeamTransfers(ctx context.Context, teamID int, opts Recursive) ([]TeamTransfer, error) {
	curPath := fmt.Sprintf("/team/%d/transfers", teamID)

	var transfers []TeamTransfer

	for {
		var resp teamTransfersResp
		if err := client.call(ctx, curPath, nil, &resp); err != nil {
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

type TeamResult struct {
	Clan1       Clan            `json:"clan1"`
	Clan2       Clan            `json:"clan2"`
	Competition CompetitionInfo `json:"competition"`
	Defaultwin  bool            `json:"defaultwin"`
	Division    Division        `json:"division"`
	Result      int             `json:"result"`
	Maps        []string        `json:"maps"`
	R1          int             `json:"r1"`
	R2          int             `json:"r2"`
	Round       string          `json:"round"`
	Time        int             `json:"time"`
	Week        int             `json:"week"`
}

type pagedTeamResponse struct {
	CurrentPage  int          `json:"current_page"`
	Data         []TeamResult `json:"data"`
	FirstPageURL string       `json:"first_page_url"`
	From         int          `json:"from"`
	LastPage     int          `json:"last_page"`
	LastPageURL  string       `json:"last_page_url"`
	Links        []links      `json:"links"`
	NextPageURL  *string      `json:"next_page_url"`
	Path         string       `json:"path"`
	PerPage      int          `json:"per_page"`
	PrevPageURL  *string      `json:"prev_page_url"`
	To           int          `json:"to"`
	Total        int          `json:"total"`
}

func (resp pagedTeamResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) TeamResults(ctx context.Context, teamID int, opts Recursive) ([]TeamResult, error) {
	curPath := fmt.Sprintf("/team/%d/results", teamID)

	var bans []TeamResult

	for {
		var resp pagedTeamResponse
		if err := client.call(ctx, curPath, nil, &resp); err != nil {
			return nil, err
		}

		bans = append(bans, resp.Data...)

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

type teamMatchesResp struct {
	CurrentPage  int          `json:"current_page"`
	Data         []TeamResult `json:"data"`
	FirstPageURL string       `json:"first_page_url"`
	From         int          `json:"from"`
	LastPage     int          `json:"last_page"`
	LastPageURL  string       `json:"last_page_url"`
	Links        []links      `json:"links"`
	NextPageURL  *string      `json:"next_page_url"`
	Path         string       `json:"path"`
	PerPage      int          `json:"per_page"`
	PrevPageURL  *string      `json:"prev_page_url"`
	To           int          `json:"to"`
	Total        int          `json:"total"`
}

func (resp teamMatchesResp) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type TeamMatchesOpts struct {
	Recursive
	// Team ID of the blu team.
	Clan1 int `json:"clan1,omitempty"`
	// Team ID of the red team.
	Clan2 int `json:"clan2,omitempty"`
	// Team ID of either team.
	Vs int `json:"vs,omitempty"`
	// If set to 1, returns matches that have yet to be played. If set to 0, returns matches that are over.
	Scheduled int `json:"scheduled,omitempty"`
	// Limit your search to a specific competition. Expects a competition ID.
	Competition int `json:"competition,omitempty"`
	// UNIX timestamp that limits results to everything after the timestamp.
	From int `json:"from,omitempty"`
	// UNIX timestamp that limits results to everything before the time.
	To int `json:"to,omitempty"`
	// Name of the division in which the competition was played.
	Division string `json:"division,omitempty"`
	// Name of the type of team.
	TeamType string `json:"team_type,omitempty"`
	// Name of the current round.
	Round string `json:"round,omitempty"`
	// A list of ETF2L user ID's. Returns only matches in which any of the provided players participated.
	Players []int `json:"players,omitempty"`
}

func (client *Client) TeamMatches(ctx context.Context, teamID int, opts Recursive) ([]TeamResult, error) {
	curPath := fmt.Sprintf("/team/%d/results", teamID)

	var bans []TeamResult

	for {
		var resp teamMatchesResp
		if err := client.call(ctx, curPath, nil, &resp); err != nil {
			return nil, err
		}

		bans = append(bans, resp.Data...)

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
