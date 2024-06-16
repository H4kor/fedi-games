package tests

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/H4kor/fedi-games/web"
	"github.com/stretchr/testify/require"
)

func TestRpsGame(t *testing.T) {
	type give struct {
		msg string
	}
	type want struct {
		contentContains string
	}
	tests := []struct {
		give
		want
	}{
		{
			give: give{"rock"},
			want: want{"I choose"},
		},
		{
			give: give{"paper"},
			want: want{"I choose"},
		},
		{
			give: give{"scissors"},
			want: want{"I choose"},
		},
		{
			give: give{"metal"},
			want: want{"Your message must contain"},
		},
		{
			give: give{""},
			want: want{"Your message must contain"},
		},
	}

	for _, test := range tests {

		server := web.DefaultServer()
		srv := server.Start()

		mock := NewMockAPServer()
		// test
		body := mock.NewNote("one", []string{"http://localhost:4040/games/rps"}, "1", "", test.give.msg)
		req, err := mock.SignedRequest("one", "POST", "/games/rps/inbox", body)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		srv.Handler.ServeHTTP(resp, req)
		require.Equal(t, resp.Result().StatusCode, 200)
		// wait for full processing
		time.Sleep(200 * time.Millisecond)

		//validation
		t.Log(mock.Retrieved)
		retrieved, ok := mock.Retrieved["one"]
		require.True(t, ok)
		require.Len(t, retrieved, 1)
		obj, _ := retrieved[0]["object"].(map[string]interface{})
		require.Contains(t, obj["content"], test.want.contentContains)

		mock.Server.Shutdown(context.Background())
	}

}
