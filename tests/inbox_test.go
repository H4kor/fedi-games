package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/H4kor/fedi-games/web"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/stretchr/testify/require"
)

func newNote(fromId string, toIds []string, noteId string, inReplyToId, msgStr string) []byte {

	to := vocab.ItemCollection{}
	mentions := vocab.ItemCollection{}
	for _, t := range toIds {
		to = append(to, vocab.ID(t))
		mention := vocab.MentionNew(
			vocab.ID(t),
		)
		mention.Href = vocab.ID(t)
		mentions = append(mentions, mention)
	}
	note := vocab.Note{
		ID:           vocab.ID("http://localhost:7777/notes/" + noteId),
		Type:         "Note",
		InReplyTo:    vocab.ID(inReplyToId),
		To:           to,
		Published:    time.Now().UTC(),
		AttributedTo: vocab.ID("http://localhost:7777/actors/" + fromId),
		Tag:          mentions,
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(msgStr)},
		},
	}

	create := vocab.CreateNew(vocab.IRI(note.ID.String()+"/activity"), note)
	create.Actor = note.AttributedTo
	create.To = note.To
	create.Published = note.Published
	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(create)

	return data

}

func TestInboxUnsigned(t *testing.T) {

	server := web.DefaultServer()
	srv := server.Start()
	// test
	req := httptest.NewRequest(
		"POST", "/games/rps/inbox",
		bytes.NewBuffer(newNote("one", []string{"http://localhost:4040/games/rps"}, "1", "", "rock")),
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
	body := newNote("one", []string{"http://localhost:4040/games/rps"}, "1", "", "rock")
	req := httptest.NewRequest(
		"POST", "/games/rps/inbox",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", "localhost:7777")
	err := Sign(mock.PrivateKey, "http://localhost:7777/actors/1#main-key", body, req)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	srv.Handler.ServeHTTP(resp, req)
	// validation
	require.Equal(t, resp.Result().StatusCode, 200)
}
