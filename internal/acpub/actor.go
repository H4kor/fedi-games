package acpub

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/H4kor/fedi-games/config"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/go-fed/httpsig"
)

func ActorToLink(act vocab.Actor) string {
	url, _ := url.Parse(act.GetLink().String())
	return "<a href=\"" + act.GetLink().String() + "\" class=\"u-url mention\">@" + act.PreferredUsername.String() + "@" + url.Host + "</a>"
}

func GetActor(reqUrl string, fromGame string) (vocab.Actor, error) {
	c := http.Client{}

	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		slog.Error("parse error", "err", err)
		return vocab.Actor{}, err
	}

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", parsedUrl.Host)

	cfg := config.GetConfig()
	err = sign(cfg.PrivKey, cfg.FullUrl()+"/games/"+fromGame+"#main-key", nil, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
		return vocab.Actor{}, err
	}

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
	headersToSign := []string{httpsig.RequestTarget, "host", "date"}
	if body != nil {
		headersToSign = append(headersToSign, "digest")
	}
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

func VerifySignature(r *http.Request, sender string, fromGame string) error {
	actor, err := GetActor(sender, fromGame)
	// actor does not have a pub key -> don't verify
	if actor.PublicKey.PublicKeyPem == "" {
		return nil
	}

	if err != nil {
		return err
	}
	block, _ := pem.Decode([]byte(actor.PublicKey.PublicKeyPem))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	slog.Info("retrieved pub key of sender", "actor", actor, "pubKey", pubKey)

	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		return err
	}
	return verifier.Verify(pubKey, httpsig.RSA_SHA256)
}

func sendObject(to vocab.Actor, fromGame string, data []byte) error {
	if to.Inbox == nil {
		slog.Error("actor has no inbox", "actor", to)
		return errors.New("actor has no inbox")
	}

	actorUrl, err := url.Parse(to.Inbox.GetID().String())
	if err != nil {
		slog.Error("parse error", "err", err)
		return err
	}

	cfg := config.GetConfig()
	c := http.Client{}
	req, _ := http.NewRequest("POST", to.Inbox.GetID().String(), bytes.NewReader(data))
	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", actorUrl.Host)
	err = sign(cfg.PrivKey, cfg.FullUrl()+"/games/"+fromGame+"#main-key", data, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
		return err
	}
	resp, err := c.Do(req)
	slog.Info("Request", "host", resp.Request.Header)

	if err != nil {
		slog.Error("Sending error", "err", err)
		return err
	}

	if resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Error sending Note", "status", resp.Status, "body", string(body))
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	slog.Info("Sent Body", "body", string(data))
	slog.Info("Retrieved", "status", resp.Status, "body", string(body))
	return nil
}

func Accept(fromGame string, act *vocab.Activity) error {
	actor, err := GetActor(act.Actor.GetID().String(), fromGame)
	if err != nil {
		return err
	}

	accept := vocab.AcceptNew(vocab.IRI("TODO"), act)
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(accept)

	if err != nil {
		slog.Error("marshalling error", "err", err)
		return err
	}

	return sendObject(actor, fromGame, data)
}

func SendNote(fromGame string, note vocab.Note) error {
	for _, to := range note.To {
		actor, err := GetActor(to.GetID().String(), fromGame)
		if err != nil {
			slog.Error("Unable to get actor", "err", err)
			return err
		}
		slog.Info("Retrieved Actor", "actor", actor, "inbox", actor.Inbox)

		create := vocab.CreateNew(vocab.IRI(note.ID.String()+"/activity"), note)
		create.Actor = note.AttributedTo
		create.To = note.To
		create.Published = note.Published
		data, err := jsonld.WithContext(
			jsonld.IRI(vocab.ActivityBaseURI),
			jsonld.Context{
				jsonld.ContextElement{
					Term: "toot",
					IRI:  jsonld.IRI("http://joinmastodon.org/ns#"),
				},
			},
		).Marshal(create)
		if err != nil {
			slog.Error("marshalling error", "err", err)
			return err
		}

		err = sendObject(actor, fromGame, data)
		if err != nil {
			return err
		}
	}
	return nil

}
