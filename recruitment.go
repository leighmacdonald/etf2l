package etf2l

import (
	"context"

	"github.com/pkg/errors"
)

type RecruitmentComments struct {
	Count int `json:"count"`
	Last  int `json:"last"`
}

type PlayerRecruitment struct {
	Classes  []string            `json:"classes"`
	Comments RecruitmentComments `json:"comments"`
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Skill    string              `json:"skill"`
	Steam    SteamPlayer         `json:"steam"`
	Type     string              `json:"type"`
	Urls     struct {
		Player      string `json:"player"`
		Recruitment string `json:"recruitment"`
	} `json:"urls"`
}

type pagedPlayerRecruitment struct {
	CurrentPage  int                 `json:"current_page"`
	Data         []PlayerRecruitment `json:"data"`
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

type playerRecruitmentResp struct {
	Pager  pagedPlayerRecruitment `json:"recruitment"`
	Status Status                 `json:"status"`
}

func (resp playerRecruitmentResp) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

type RecruitmentOpts struct {
	BaseOpts
	// Returns only recruitment posts of a specific country.
	Country string `json:"country,omitempty"`
	// Returns only recruitment posts of a specific class. Can be provided as string or as a list.
	// In order to search for multiple classes, provide the argument in an array/list format.
	Class []string `json:"class,omitempty"`
	// Returns only recruitment posts for a certain skill level. Can be provided as string or as a list.
	// In order to search for multiple skill levels, provide the argument in an array/list format.
	Skill []string `json:"skill,omitempty"`
	// Limit recruitment posts by team type.
	Type string `json:"type,omitempty"`
	// Limit recruitment posts by ETF2L user id. Is the creator of the post.
	User int `json:"user,omitempty"`
}

func (client *Client) PlayerRecruitment(ctx context.Context, opts RecruitmentOpts) ([]PlayerRecruitment, error) {
	var matches []PlayerRecruitment

	curPath := "/recruitment/players"

	for {
		var resp playerRecruitmentResp
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

type TeamRecruitment struct {
	Classes  []string            `json:"classes"`
	Comments RecruitmentComments `json:"comments"`
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Skill    string              `json:"skill"`
	Steam    SteamPlayer         `json:"steam"`
	Type     string              `json:"type"`
	Urls     struct {
		Team        string `json:"team"`
		Recruitment string `json:"recruitment"`
	} `json:"urls"`
}

type pagedTeamRecruitment struct {
	CurrentPage  int               `json:"current_page"`
	Data         []TeamRecruitment `json:"data"`
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

type teamRecruitmentResp struct {
	Pager  pagedTeamRecruitment `json:"recruitment"`
	Status Status               `json:"status"`
}

func (resp teamRecruitmentResp) NextURL(r Recursive) (string, error) {
	if !r.IsRecursive() || resp.Pager.NextPageURL == nil {
		return "", ErrEOF
	}

	nextPath, err := getPath(*resp.Pager.NextPageURL)
	if err != nil {
		return "", err
	}

	return nextPath, nil
}

func (client *Client) TeamRecruitment(ctx context.Context, opts RecruitmentOpts) ([]TeamRecruitment, error) {
	var matches []TeamRecruitment

	curPath := "/recruitment/teams"

	for {
		var resp teamRecruitmentResp
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
