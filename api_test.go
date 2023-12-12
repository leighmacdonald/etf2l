package etf2l_test

import (
	"context"
	"github.com/leighmacdonald/etf2l"
	"github.com/leighmacdonald/steamid/v3/steamid"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testIDb4nny  steamid.SID64 = "76561197970669109"
	testIDBanned steamid.SID64 = "76561198203516436"
)

func TestPlayer(t *testing.T) {
	client := etf2l.New()
	p1, err := client.Player(context.Background(), testIDBanned.String())
	require.NoError(t, err)
	require.Equal(t, 139491, p1.Id)

	_, err404 := client.Player(context.Background(), "7999198203516436")
	require.ErrorIs(t, etf2l.ErrNotFound, err404)
}

func TestWhitelists(t *testing.T) {
	client := etf2l.New()
	whitelists, err := client.Whitelists(context.Background())
	require.NoError(t, err)
	require.True(t, len(whitelists) > 4)
}

func TestDemos(t *testing.T) {
	client := etf2l.New()
	demos, err := client.Demos(context.Background(), etf2l.DemoOpts{PlayerID: "2788"})
	require.NoError(t, err)
	require.True(t, len(demos) > 20)
}
