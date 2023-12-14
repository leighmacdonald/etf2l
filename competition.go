package etf2l

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Competition struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Archived    bool   `json:"archived"`
	Type        string `json:"type"`
	Urls        struct {
		Matches string `json:"matches"`
		Results string `json:"results"`
		Self    string `json:"self"`
		Teams   string `json:"teams"`
	} `json:"urls"`
}

type pagedCompetitions struct {
	CurrentPage  int           `json:"current_page"`
	Data         []Competition `json:"data"`
	FirstPageURL string        `json:"first_page_url"`
	From         int           `json:"from"`
	LastPage     int           `json:"last_page"`
	LastPageURL  string        `json:"last_page_url"`
	Links        []links       `json:"links"`
	NextPageURL  *string       `json:"next_page_url"`
	Path         string        `json:"path"`
	PerPage      int           `json:"per_page"`
	PrevPageURL  *string       `json:"prev_page_url"`
	To           int           `json:"to"`
	Total        int           `json:"total"`
}

type competitionResponse struct {
	Pager pagedCompetitions `json:"competitions"`
}

func (resp competitionResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type ArchivedState int

const (
	Active = iota
	Archived
)

type CompetitionOpts struct {
	Recursive
	Archived    ArchivedState `json:"archived,omitempty"` // 1 Returns only archived competitions, 0 returns non-archived competitions.
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Category    string        `json:"category,omitempty"`
	CompType    string        `json:"comp_type,omitempty"`
	TeamType    string        `json:"team_type,omitempty"`
	Competition string        `json:"competition,omitempty"`
}

func (client *Client) CompetitionList(ctx context.Context, opts Recursive) ([]Competition, error) {
	var competitions []Competition

	curPath := "/competition/list"

	for {
		var resp competitionResponse
		if err := client.call(ctx, curPath, opts, &resp); err != nil {
			return nil, err
		}

		competitions = append(competitions, resp.Pager.Data...)

		nextURL, err := resp.NextURL(opts)
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return competitions, nil
}

type CompetitionDetails struct {
	Category    string   `json:"category"`
	Description string   `json:"description"`
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Pool        []string `json:"pool"`
	Archived    bool     `json:"archived"`
	Type        string   `json:"type"`
	Teams       struct {
		Max      int `json:"max"`
		Signedup int `json:"signedup"`
	} `json:"teams"`
	Urls struct {
		Matches string `json:"matches"`
		Results string `json:"results"`
		Self    string `json:"self"`
		Teams   string `json:"teams"`
	} `json:"urls"`
}

type competitionDetailsResponse struct {
	Competition CompetitionDetails `json:"competition"`
	Status      Status             `json:"status"`
}

func (client *Client) CompetitionDetails(ctx context.Context, competitionID int) (CompetitionDetails, error) {
	var resp competitionDetailsResponse
	if err := client.call(ctx, fmt.Sprintf("/competition/%d", competitionID), nil, &resp); err != nil {
		return CompetitionDetails{}, err
	}

	return resp.Competition, nil
}

type CompetitionTeam struct {
	ID      int        `json:"id"`
	Country string     `json:"country"`
	Name    string     `json:"name"`
	Dropped int        `json:"dropped"`
	Steam   SteamGroup `json:"steam"`
	URL     string     `json:"url"`
}

type pagedCompetitionTeams struct {
	CurrentPage  int               `json:"current_page"`
	Data         []CompetitionTeam `json:"data"`
	FirstPageURL string            `json:"first_page_url"`
	From         int               `json:"from"`
	LastPage     int               `json:"last_page"`
	LastPageURL  string            `json:"last_page_url"`
	Links        []links           `json:"links"`
	NextPageURL  *string           `json:"next_page_url"`
	Path         string            `json:"path"`
	PerPage      int               `json:"per_page"`
	PrevPageURL  *string           `json:"prev_page_url"`
	To           int               `json:"to"`
	Total        int               `json:"total"`
}

type competitionTeamsResponse struct {
	Pager  pagedCompetitionTeams `json:"teams"`
	Status Status                `json:"status"`
}

func (resp competitionTeamsResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) CompetitionTeams(ctx context.Context, competitionID int, opts BaseOpts) ([]CompetitionTeam, error) {
	var teams []CompetitionTeam

	curPath := fmt.Sprintf("/competition/%d/teams", competitionID)

	for {
		var resp competitionTeamsResponse
		if err := client.call(ctx, curPath, nil, &resp); err != nil {
			return nil, err
		}

		teams = append(teams, resp.Pager.Data...)

		nextURL, err := resp.NextURL(opts)
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return nil, err
		}

		curPath = nextURL
	}

	return teams, nil
}

type Clan struct {
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

type CompetitionInfo struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type CompetitionResult struct {
	Clan1       Clan            `json:"clan1"`
	Clan2       Clan            `json:"clan2"`
	Competition CompetitionInfo `json:"competition"`
	Defaultwin  bool            `json:"defaultwin"`
	Division    Division        `json:"division"`
	ID          int             `json:"id"`
	Maps        []string        `json:"maps"`
	R1          int             `json:"r1"`
	R2          int             `json:"r2"`
	Round       string          `json:"round"`
	Time        interface{}     `json:"time"`
	Week        int             `json:"week"`
}

type pagedCompetitionResults struct {
	CurrentPage  int                 `json:"current_page"`
	Data         []CompetitionResult `json:"data"`
	FirstPageURL string              `json:"first_page_url"`
	From         int                 `json:"from"`
	LastPage     int                 `json:"last_page"`
	LastPageURL  string              `json:"last_page_url"`
	Links        []links             `json:"links"`
	NextPageURL  *string             `json:"next_page_url"`
	Path         string              `json:"path"`
	PerPage      int                 `json:"per_page"`
	PrevPageURL  *string             `json:"prev_page_url"`
	To           int                 `json:"to"`
	Total        int                 `json:"total"`
}

type competitionResultsResponse struct {
	Pager  pagedCompetitionResults `json:"results"`
	Status Status                  `json:"status"`
}

func (resp competitionResultsResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) CompetitionResults(ctx context.Context, competitionID int, opts Recursive) ([]CompetitionResult, error) {
	var results []CompetitionResult

	curPath := fmt.Sprintf("/competition/%d/results", competitionID)

	for {
		var resp competitionResultsResponse
		if err := client.call(ctx, curPath, nil, &resp); err != nil {
			return nil, err
		}

		results = append(results, resp.Pager.Data...)

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

type CompetitionMatchClan struct {
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

type CompetitionMatch struct {
	Clan1       CompetitionMatchClan `json:"clan1"`
	Clan2       CompetitionMatchClan `json:"clan2"`
	Competition struct {
		Category string `json:"category"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		URL      string `json:"url"`
	} `json:"competition"`
	Defaultwin bool `json:"defaultwin"`
	Division   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Tier int    `json:"tier"`
	} `json:"division"`
	ID     int      `json:"id"`
	Maps   []string `json:"maps"`
	Result struct {
		R1 int `json:"r1"`
		R2 int `json:"r2"`
	} `json:"result"`
	Round        string      `json:"round"`
	Time         interface{} `json:"time"`
	Week         int         `json:"week"`
	SkillContrib int         `json:"skill_contrib"`
}

type pagedCompetitionMatches struct {
	CurrentPage  int                `json:"current_page"`
	Data         []CompetitionMatch `json:"data"`
	FirstPageURL string             `json:"first_page_url"`
	From         int                `json:"from"`
	LastPage     int                `json:"last_page"`
	LastPageURL  string             `json:"last_page_url"`
	Links        []links            `json:"links"`
	NextPageURL  *string            `json:"next_page_url"`
	Path         string             `json:"path"`
	PerPage      int                `json:"per_page"`
	PrevPageURL  *string            `json:"prev_page_url"`
	To           int                `json:"to"`
	Total        int                `json:"total"`
}

type competitionMatchesResponse struct {
	Pager  pagedCompetitionMatches `json:"matches"`
	Status Status                  `json:"status"`
}

func (resp competitionMatchesResponse) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) CompetitionMatches(ctx context.Context, competitionID int, opts Recursive) ([]CompetitionMatch, error) {
	var matches []CompetitionMatch

	curPath := fmt.Sprintf("/competition/%d/matches", competitionID)

	for {
		var resp competitionMatchesResponse
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

type CompetitionTable struct {
	ID            int    `json:"id"`
	Drop          bool   `json:"drop"`
	DivisionID    int    `json:"division_id"`
	DivisionName  string `json:"division_name"`
	Country       string `json:"country"`
	Name          string `json:"name"`
	MapsPlayed    int    `json:"maps_played"`
	MapsWon       int    `json:"maps_won"`
	GcWon         int    `json:"gc_won"`
	GcLost        int    `json:"gc_lost"`
	MapsLost      int    `json:"maps_lost"`
	PenaltyPoints int    `json:"penalty_points"`
	Score         int    `json:"score"`
	Ach           int    `json:"ach"`
	Byes          int    `json:"byes"`
	SeededPoints  int    `json:"seeded_points"`
}

type TablesResponse struct {
	Tables map[string]CompetitionTable
}

func (client *Client) CompetitionTables(ctx context.Context, competitionID int) (map[string]CompetitionTable, error) {
	var resp TablesResponse
	if err := client.call(ctx, fmt.Sprintf("/competition/%d/tables", competitionID), nil, &resp); err != nil {
		return nil, err
	}

	return resp.Tables, nil
}
