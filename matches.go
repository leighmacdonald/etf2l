package etf2l

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type MatchClan struct {
	Country string `json:"country"`
	Drop    bool   `json:"drop"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Steam   struct {
		Avatar string `json:"avatar"`
		Group  string `json:"group"`
	} `json:"steam"`
	URL string `json:"url"`
}

type MatchCompetition struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type Match struct {
	Clan1       MatchClan        `json:"clan1"`
	Clan2       MatchClan        `json:"clan2"`
	Competition MatchCompetition `json:"competition"`
	Submitted   int              `json:"submitted"`
	Defaultwin  bool             `json:"defaultwin"`
	Division    Division         `json:"division"`
	ID          int              `json:"id"`
	Maps        []string         `json:"maps"`
	R1          int              `json:"r1"`
	R2          int              `json:"r2"`
	Round       string           `json:"round"`
	Time        int              `json:"time"`
	Week        int              `json:"week"`
	Urls        struct {
		Self string `json:"self"`
		API  string `json:"api"`
	} `json:"urls"`
}

type pagedMatches struct {
	CurrentPage  int     `json:"current_page"`
	Data         []Match `json:"data"`
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

type matchesResponse struct {
	Pager  pagedMatches `json:"results"`
	Status Status       `json:"status"`
}

func (resp matchesResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type MatchesOpts struct {
	BaseOpts
	Clan1       int      `url:"clan1,omitempty"`       // Team ID of the blu team.
	Clan2       int      `url:"clan2,omitempty"`       // Team ID of the red team.
	Vs          int      `url:"vs,omitempty"`          // Team ID of either team.
	Scheduled   int      `url:"scheduled,omitempty"`   // If set to 1, returns matches that have yet to be played. If set to 0, returns matches that are over.
	Competition int      `url:"competition,omitempty"` // Limit your search to a specific competition. Expects a competition ID.
	From        int      `url:"from,omitempty"`        // UNIX timestamp that limits results to everything after the timestamp.
	To          int      `url:"to,omitempty"`          // UNIX timestamp that limits results to everything before the time.
	Division    string   `url:"division,omitempty"`    // Name of the division in which the competition was played.
	TeamType    string   `url:"team_type,omitempty"`   // Name of the type of team.
	Round       string   `url:"round,omitempty"`       // Name of the current round.
	Players     []string `url:"players,omitempty"`     // A list of ETF2L user ID's. Returns only matches in which any of the provided players participated.
}

func (client *Client) Matches(ctx context.Context, opts Recursive) ([]Match, error) {
	var matches []Match

	curPath := "/matches"

	for {
		var resp matchesResponse
		if err := client.call(ctx, curPath, nil, &resp); err != nil {
			return nil, err
		}

		matches = append(matches, resp.Pager.Data...)

		nextURL, err := resp.NextURL(opts)
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return matches, nil
}

type MatchMapResult struct {
	MatchOrder int    `json:"match_order"`
	Clan1      int    `json:"clan1"`
	Clan2      int    `json:"clan2"`
	Map        string `json:"map"`
	GoldenCap  bool   `json:"golden_cap"`
}

type MatchDetails struct {
	Clan1       MatchClan        `json:"clan1"`
	Clan2       MatchClan        `json:"clan2"`
	Competition MatchCompetition `json:"competition"`
	Defaultwin  bool             `json:"defaultwin"`
	Division    interface{}      `json:"division"`
	ID          int              `json:"id"`
	Maps        []string         `json:"maps"`
	R1          int              `json:"r1"`
	R2          int              `json:"r2"`
	Round       string           `json:"round"`
	Time        int              `json:"time"`
	Submitted   int              `json:"submitted"`
	Week        int              `json:"week"`
	Urls        struct {
		Self string `json:"self"`
		API  string `json:"api"`
	} `json:"urls"`
	Players    []interface{}    `json:"players"`
	ByeWeek    bool             `json:"bye_week"`
	Demos      []interface{}    `json:"demos"`
	MapResults []MatchMapResult `json:"map_results"`
}

type matchDetailsResponse struct {
	Match  MatchDetails `json:"match"`
	Status Status       `json:"status"`
}

func (client *Client) MatchDetails(ctx context.Context, leagueMatchID int) (*MatchDetails, error) {
	var resp matchDetailsResponse
	if err := client.call(ctx, fmt.Sprintf("/matches/%d", leagueMatchID), nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Match, nil
}
