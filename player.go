package etf2l

import (
	"context"
	"fmt"
	"net/http"
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

type Division struct {
	Name         string `json:"name"`
	SkillContrib int    `json:"skill_contrib"`
	Tier         int    `json:"tier"`
}

type Competition struct {
	Category    string   `json:"category"`
	Competition string   `json:"competition"`
	Division    Division `json:"division"`
	Url         string   `json:"url"`
}

type IRC struct {
	Channel *string `json:"channel"`
	Network *string `json:"network"`
}

type SteamGroup struct {
	Avatar     string  `json:"avatar"`
	SteamGroup *string `json:"steam_group"`
}

type SteamPlayer struct {
	Avatar string `json:"avatar"`
	Id     string `json:"id"`
	Id3    string `json:"id3"`
	Id64   string `json:"id64"`
}

type Team struct {
	Competitions map[string]Competition `json:"competitions"`
	Country      string                 `json:"country"`
	Homepage     *string                `json:"homepage"`
	Id           int                    `json:"id"`
	Irc          IRC                    `json:"irc"`
	Name         string                 `json:"name"`
	Server       *string                `json:"server"`
	Steam        SteamGroup             `json:"steam"`
	Tag          string                 `json:"tag"`
	Type         string                 `json:"type"`
	Urls         URLs                   `json:"urls"`
}

type URLs struct {
	Matches   string `json:"matches"`
	Results   string `json:"results"`
	Self      string `json:"self"`
	Transfers string `json:"transfers"`
}

type Player struct {
	Bans       []BanReason `json:"bans"`
	Classes    []string    `json:"classes"`
	Country    string      `json:"country"`
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Registered int         `json:"registered"`
	Steam      SteamPlayer `json:"steam"`
	Teams      []Team      `json:"teams"`
	Title      string      `json:"title"`
	Urls       struct {
		Results   string `json:"results"`
		Self      string `json:"self"`
		Transfers string `json:"transfers"`
	} `json:"urls"`
}

func (client *Client) Player(ctx context.Context, playerID string) (*Player, error) {
	var resp PlayerResponse
	if err := client.call(ctx, http.MethodGet, fmt.Sprintf("/player/%s", playerID), nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Player, nil
}
