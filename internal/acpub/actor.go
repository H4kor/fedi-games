package acpub

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/go-fed/httpsig"
	"rerere.org/fedi-games/config"
)

func GetActor(url string) (vocab.Actor, error) {
	c := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/ld+json")

	resp, err := c.Do(req)
	if err != nil {
		return vocab.Actor{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return vocab.Actor{}, err
	}

	item, err := vocab.UnmarshalJSON(data)
	if err != nil {
		return vocab.Actor{}, err
	}

	var actor vocab.Actor

	err = vocab.OnActor(item, func(o *vocab.Actor) error {
		actor = *o
		return nil
	})

	return actor, err
}

func sign(privateKey *rsa.PrivateKey, pubKeyId string, body []byte, r *http.Request) error {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256
	// The "Date" and "Digest" headers must already be set on r, as well as r.URL.
	headersToSign := []string{httpsig.RequestTarget, "host", "date", "digest"}
	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, 0)
	if err != nil {
		return err
	}
	// To sign the digest, we need to give the signer a copy of the body...
	// ...but it is optional, no digest will be signed if given "nil"
	// If r were a http.ResponseWriter, call SignResponse instead.
	err = signer.SignRequest(privateKey, pubKeyId, r, body)

	slog.Info("Signed Request", "req", r.Header)
	return err
}

func SendNote(fromGame string, note vocab.Note) error {
	cfg := config.GetConfig()

	actor, err := GetActor(note.To[0].GetID().String())
	if err != nil {
		slog.Error("Unable to get actor", "err", err)
		return err
	}
	slog.Info("Retrieved Actor", "actor", actor, "inbox", actor.Inbox)

	create := vocab.CreateNew(vocab.IRI(note.ID.String()+"/activity"), note)
	create.Actor = note.AttributedTo
	create.To = note.To
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(create)
	if err != nil {
		return err
	}

	actorUrl, err := url.Parse(actor.Inbox.GetID().String())
	if err != nil {
		return err
	}

	c := http.Client{}
	req, _ := http.NewRequest("POST", actor.Inbox.GetID().String(), bytes.NewReader(data))
	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", actorUrl.Host)
	err = sign(cfg.PrivKey, cfg.FullUrl()+"/games/"+fromGame+"#main-key", data, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
	}
	resp, err := c.Do(req)

	slog.Info("Request", "host", resp.Request.Header)

	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Error sending Note", "status", resp.Status, "body", string(body))
		return errors.New("error status code " + string(resp.Status))
	}
	body, _ := io.ReadAll(resp.Body)
	slog.Info("Sent Body", "body", string(data))
	slog.Info("Retrieved", "status", resp.Status, "body", string(body))
	return nil
}