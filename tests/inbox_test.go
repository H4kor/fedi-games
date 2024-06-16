package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/H4kor/fedi-games/web"
	"github.com/stretchr/testify/require"
)

func TestInboxUnsigned(t *testing.T) {

	server := web.DefaultServer()
	srv := server.Start()
	mock := NewMockAPServer()
	defer mock.Server.Shutdown(context.Background())
	// test
	req := httptest.NewRequest(
		"POST", "/games/rps/inbox",
		bytes.NewBuffer(mock.NewNote("one", []string{"http://localhost:4040/games/rps"}, "1", "", "rock")),
	)
	resp := httptest.NewRecorder()
	srv.Handler.ServeHTTP(resp, req)
	// validation
	require.Equal(t, resp.Result().StatusCode, 401)

}

func TestInboxSigned(t *testing.T) {

	server := web.DefaultServer()
	srv := server.Start()

	mock := NewMockAPServer()
	defer mock.Server.Shutdown(context.Background())
	// test
	body := mock.NewNote("one", []string{"http://localhost:4040/games/rps"}, "1", "", "rock")
	req := httptest.NewRequest(
		"POST", "/games/rps/inbox",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", "localhost:7777")
	err := Sign(mock.PrivateKey, "http://localhost:7777/actors/one#main-key", body, req)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	srv.Handler.ServeHTTP(resp, req)
	// validation
	require.Equal(t, resp.Result().StatusCode, 200)
}
