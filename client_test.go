package etf2l_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/leighmacdonald/etf2l"
	"github.com/leighmacdonald/steamid/v4/steamid"
	"github.com/stretchr/testify/require"
)

const (
	testETF2LBannedID int = 139491
)

var (
	testIDb4nny  = steamid.New("76561197970669109")
	testIDBanned = steamid.New("76561198203516436")
)

func TestClient(t *testing.T) {
	client := etf2l.New()

	t.Run("player", testPlayer(client))
	t.Run("player_results", testPlayerResults(client))
	t.Run("player_transfers", testPlayerTransfers(client))
	t.Run("demos", testDemos(client))
	t.Run("bans", testBans(client))
	t.Run("bans_recursive", testBansRecursive(client))
	t.Run("competition_list", testCompetitionList(client))
	t.Run("competition_details", testCompetitionDetails(client))
	t.Run("competition_teams", testCompetitionTeams(client))
	t.Run("competition_results", testCompetitionResults(client))
	t.Run("competition_matches", testCompetitionMatches(client))
	t.Run("competition_tables", testCompetitionTables(client))
	t.Run("matches", testMatches(client))
	t.Run("match_details", testMatchDetails(client))
	t.Run("whitelist", testWhitelists(client))
	t.Run("player_recruitment", testPlayerRecruitment(client))
	t.Run("team_recruitment", testTeamRecruitment(client))
	t.Run("team", testTeam(client))
	t.Run("team_transfers", testTeamTransfers(client))
	t.Run("team_results", testTeamResults(client))
	t.Run("team_matches", testTeamMatches(client))
}

func testPlayer(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		p1, err := client.Player(context.Background(), &http.Client{}, testIDb4nny.String())
		require.NoError(t, err)
		require.Equal(t, 20834, p1.ID)

		_, err404 := client.Player(context.Background(), &http.Client{}, "7999198203516436")
		require.ErrorIs(t, etf2l.ErrNotFound, err404)
	}
}

func testPlayerResults(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		results, err := client.PlayerResults(context.Background(), &http.Client{}, testIDBanned.String(), etf2l.BaseOpts{Recursive: false})
		require.NoError(t, err)
		require.Equal(t, 20, len(results))
	}
}

func testPlayerTransfers(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		results, err := client.PlayerTransfers(context.Background(), &http.Client{}, testETF2LBannedID, etf2l.BaseOpts{Recursive: false})
		require.NoError(t, err)
		require.Equal(t, 20, len(results))
	}
}

func testDemos(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		demos, err := client.Demos(context.Background(), &http.Client{}, etf2l.DemoOpts{
			Recursive: etf2l.BaseOpts{Recursive: false},
			PlayerID:  "2788",
		})
		require.NoError(t, err)
		require.Equal(t, 20, len(demos))
	}
}

func testBans(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		bans, err := client.Bans(context.Background(), &http.Client{}, etf2l.BanOpts{
			Recursive: etf2l.BaseOpts{Recursive: false},
			PlayerID:  testETF2LBannedID,
		})
		require.NoError(t, err)
		require.True(t, len(bans) > 2)
	}
}

func testBansRecursive(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		bans, err := client.Bans(context.Background(), &http.Client{}, etf2l.BanOpts{
			Recursive: etf2l.BaseOpts{Recursive: true},
		})
		require.NoError(t, err)
		require.True(t, len(bans) > 3000)
	}
}

func testCompetitionList(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		competitions, err := client.CompetitionList(context.Background(), &http.Client{}, etf2l.CompetitionOpts{
			Recursive: etf2l.BaseOpts{Recursive: false},
		})
		require.NoError(t, err)
		require.Equal(t, 20, len(competitions))
	}
}

func testCompetitionDetails(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		competition, err := client.CompetitionDetails(context.Background(), &http.Client{}, 1)
		require.NoError(t, err)
		require.Equal(t, 1, competition.ID)
	}
}

func testCompetitionTeams(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		teams, err := client.CompetitionTeams(context.Background(), &http.Client{}, 1, etf2l.BaseOpts{
			Recursive: false,
		})
		require.NoError(t, err)
		require.True(t, len(teams) == 20)
	}
}

func testCompetitionResults(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		teams, err := client.CompetitionResults(context.Background(), &http.Client{}, 1, etf2l.BaseOpts{
			Recursive: false,
		})
		require.NoError(t, err)
		require.True(t, len(teams) == 20)
	}
}

func testCompetitionMatches(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		teams, err := client.CompetitionMatches(context.Background(), &http.Client{}, 1, etf2l.BaseOpts{
			Recursive: false,
		})
		require.NoError(t, err)
		require.True(t, len(teams) == 20)
	}
}

func testCompetitionTables(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		teams, err := client.CompetitionResults(context.Background(), &http.Client{}, 1, etf2l.BaseOpts{
			Recursive: false,
		})
		require.NoError(t, err)
		require.True(t, len(teams) > 5)
	}
}

func testMatches(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		pagesData, err := client.MatchesPage(context.Background(), &http.Client{}, 0, 2000)
		require.NoError(t, err)
		require.True(t, len(pagesData.Pager.Data) > 5)
		require.True(t, pagesData.Pager.Total > 5)
	}
}

func testMatchDetails(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		match, err := client.MatchDetails(context.Background(), &http.Client{}, 1)
		require.NoError(t, err)
		require.Equal(t, 1, match.ID)
	}
}

func testWhitelists(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		whitelists, err := client.Whitelists(context.Background(), &http.Client{})
		require.NoError(t, err)
		require.True(t, len(whitelists) > 4)
	}
}

func testPlayerRecruitment(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		recruitments, err := client.PlayerRecruitment(context.Background(), &http.Client{}, etf2l.RecruitmentOpts{
			BaseOpts: etf2l.BaseOpts{Recursive: false},
		})
		require.NoError(t, err)
		require.Equal(t, 20, len(recruitments))
	}
}

func testTeamRecruitment(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		recruitments, err := client.TeamRecruitment(context.Background(), &http.Client{}, etf2l.RecruitmentOpts{
			BaseOpts: etf2l.BaseOpts{Recursive: false},
		})
		require.NoError(t, err)
		require.Equal(t, 20, len(recruitments))
	}
}

func testTeam(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		recruitments, err := client.Team(context.Background(), &http.Client{}, 2)
		require.NoError(t, err)
		require.Greater(t, len(recruitments.Competitions), 10)
	}
}

func testTeamTransfers(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		transfers, err := client.TeamTransfers(context.Background(), &http.Client{}, 2, etf2l.BaseOpts{Recursive: false})
		require.NoError(t, err)
		require.Equal(t, 20, len(transfers))
	}
}

func testTeamResults(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		results, err := client.TeamResults(context.Background(), &http.Client{}, 2, etf2l.BaseOpts{Recursive: false})
		require.NoError(t, err)
		require.Equal(t, 20, len(results))
	}
}

func testTeamMatches(client *etf2l.Client) func(*testing.T) {
	return func(t *testing.T) {
		results, err := client.TeamMatches(context.Background(), &http.Client{}, 2, etf2l.BaseOpts{Recursive: false})
		require.NoError(t, err)
		require.Equal(t, 20, len(results))
	}
}
