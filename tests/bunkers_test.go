package tests

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/H4kor/fedi-games/web"
	"github.com/stretchr/testify/require"
)

func TestBunkersGame(t *testing.T) {
	server := web.DefaultServer()
	srv := server.Start()
	mock := NewMockAPServer()

	// first round
	{
		body := mock.NewNote(
			"one", []string{"http://localhost:4040/games/bunkers", mock.ActorUrl("two")},
			"1", "", "angle 45 power 50",
		)
		req, err := mock.SignedRequest("one", "POST", "/games/bunkers/inbox", body)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		srv.Handler.ServeHTTP(resp, req)
		require.Equal(t, resp.Result().StatusCode, 200)
		// wait for full processing
		time.Sleep(200 * time.Millisecond)
		//validation
		// reply to one
		retrieved, ok := mock.Retrieved["one"]
		require.True(t, ok)
		require.Len(t, retrieved, 1)
		// reply to two
		{
			retrieved, ok := mock.Retrieved["one"]
			require.True(t, ok)
			require.Len(t, retrieved, 1)
		}
		obj, _ := retrieved[0]["object"].(map[string]interface{})
		require.NotNil(t, obj["attachment"])
	}

	// second round
	{
		body := mock.NewNote(
			"two", []string{"http://localhost:4040/games/bunkers", mock.ActorUrl("one")},
			"1", "", "angle 45 power 50",
		)
		req, err := mock.SignedRequest("one", "POST", "/games/bunkers/inbox", body)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		srv.Handler.ServeHTTP(resp, req)
		require.Equal(t, resp.Result().StatusCode, 200)
		// wait for full processing
		time.Sleep(200 * time.Millisecond)
		//validation
		retrieved, ok := mock.Retrieved["two"]
		require.True(t, ok)
		require.Len(t, retrieved, 2)
		// reply to one
		{
			retrieved, ok := mock.Retrieved["two"]
			require.True(t, ok)
			require.Len(t, retrieved, 2)
		}
		obj, _ := retrieved[1]["object"].(map[string]interface{})
		require.NotNil(t, obj["attachment"])
	}

}

func TestBunkersGameEnd(t *testing.T) {
	server := web.DefaultServer()
	srv := server.Start()
	mock := NewMockAPServer()

	// first round
	body := mock.NewNote(
		"one", []string{"http://localhost:4040/games/bunkers", mock.ActorUrl("two")},
		"1", "", "angle 90 power 1",
	)
	req, err := mock.SignedRequest("one", "POST", "/games/bunkers/inbox", body)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	srv.Handler.ServeHTTP(resp, req)
	require.Equal(t, resp.Result().StatusCode, 200)
	// wait for full processing
	time.Sleep(200 * time.Millisecond)
	//validation
	// reply to one
	retrieved, ok := mock.Retrieved["one"]
	require.True(t, ok)
	require.Len(t, retrieved, 1)
	// reply to two
	{
		retrieved, ok := mock.Retrieved["one"]
		require.True(t, ok)
		require.Len(t, retrieved, 1)
	}
	obj, _ := retrieved[0]["object"].(map[string]interface{})
	require.NotNil(t, obj["attachment"])
	attachment := obj["attachment"].(map[string]interface{})
	// gif on end of round
	require.Equal(t, "image/gif", attachment["mediaType"])

}
