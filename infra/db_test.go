package infra_test

import (
	"testing"

	"github.com/H4kor/fedi-games/games"
	"github.com/H4kor/fedi-games/infra"
	"github.com/stretchr/testify/require"
)

func TestReplySave(t *testing.T) {
	r := games.GameReply{
		To:  []string{"a", "b"},
		Msg: "msg",
		Attachments: []games.GameAttachment{
			{Url: "http://example.com", MediaType: "image/png"},
		},
	}

	db := infra.NewDatabase(":memory:")

	err := db.PersistGameReply("game", &r)
	require.NoError(t, err)
	require.NotZero(t, r.Id)

	rr, err := db.RetrieveGameReply("game", r.Id)
	require.NoError(t, err)

	require.Equal(t, r.Id, rr.Id)
}
