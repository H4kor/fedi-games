package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/H4kor/fedi-games/games/rps"
	"github.com/H4kor/fedi-games/web"
	"github.com/stretchr/testify/require"
)

func TestGetActor(t *testing.T) {
	type give struct {
		path   string
		accept string
	}
	type want struct {
		status  int
		contain string
	}
	tests := []struct {
		give
		want
	}{
		{
			give: give{"/games/rps", "application/activity+json"},
			want: want{200, rps.NewRockPaperScissorGame().Summary()},
		},
		{
			give: give{"/games/rps", "text/html"},
			want: want{200, rps.NewRockPaperScissorGame().Summary()},
		},
		{
			give: give{"/games/no-game", "application/activity+json"},
			want: want{404, ""},
		},
		{
			give: give{"/games/no-game", "text/html"},
			want: want{404, ""},
		},
	}

	server := web.DefaultServer()
	srv := server.Start()
	for _, test := range tests {
		// test
		req := httptest.NewRequest("GET", test.give.path, nil)
		req.Header.Set("Accept", test.give.accept)
		resp := httptest.NewRecorder()
		srv.Handler.ServeHTTP(resp, req)
		// validation
		require.Equal(t, resp.Result().StatusCode, test.want.status)
		require.Contains(t, resp.Body.String(), test.want.contain)
	}

}
