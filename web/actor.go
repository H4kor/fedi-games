package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/games"
)

type htmlData struct {
	Game games.Game
	Cfg  config.Config
}

type PropertyValue struct {
	Type  vocab.ActivityVocabularyType `jsonld:"type,omitempty"`
	Name  vocab.NaturalLanguageValues  `jsonld:"name"`
	Value vocab.NaturalLanguageValues  `jsonld:"value"`
}

// GetID implements activitypub.ObjectOrLink.
func (p PropertyValue) GetID() vocab.IRI {
	println("GetID")
	return "123"
}

// GetLink implements activitypub.ObjectOrLink.
func (p PropertyValue) GetLink() vocab.IRI {
	println("GetLink")
	return "123"
}

// GetType implements activitypub.ObjectOrLink.
func (p PropertyValue) GetType() vocab.ActivityVocabularyType {
	println("GetType")
	return p.Type
}

// IsCollection implements activitypub.ObjectOrLink.
func (p PropertyValue) IsCollection() bool {
	println("IsCollection")
	return false
}

// IsLink implements activitypub.ObjectOrLink.
func (p PropertyValue) IsLink() bool {
	println("IsLink")
	return false
}

// IsObject implements activitypub.ObjectOrLink.
func (p PropertyValue) IsObject() bool {
	println("IsObject")
	return true
}

// MarshalJSON encodes the receiver object to a JSON document.
func (o PropertyValue) MarshalJSON() ([]byte, error) {
	type oJson struct {
		Type  string
		Name  string
		Value string
	}

	data := oJson{
		Type:  string(o.Type),
		Name:  o.Name.First().Value.String(),
		Value: o.Value.First().Value.String(),
	}

	return json.Marshal(data)
}

func GameHandlerHtml(w http.ResponseWriter, r *http.Request, game games.Game) {
	cfg := config.GetConfig()
	err := RenderTemplateWithBase(w, "game", htmlData{
		Game: game,
		Cfg:  cfg,
	})
	if err != nil {
		slog.Error("Error rendering index", "err", err)
	}
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	game, ok := GameMap[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		GameHandlerHtml(w, r, game)
		return
	}

	cfg := config.GetConfig()

	actor := vocab.ServiceNew(vocab.IRI(cfg.FullUrl() + "/games/" + game.Name()))
	actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(game.Name())}}
	actor.Inbox = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/inbox")
	actor.Following = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/following")
	actor.Outbox = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/outbox")
	// actor.Followers = vocab.IRI(config.FullUrl() + "/games/" + game + "/followers")
	actor.PublicKey = vocab.PublicKey{
		ID:           vocab.ID(cfg.FullUrl() + "/games/" + game.Name() + "#main-key"),
		Owner:        vocab.IRI(cfg.FullUrl() + "/games/" + game.Name()),
		PublicKeyPem: cfg.PublicKeyPem,
	}

	actor.Name = vocab.NaturalLanguageValues{{Value: vocab.Content(game.Name())}}
	actor.Icon = vocab.Image{
		Type:      vocab.ImageType,
		MediaType: "image/png",
		URL:       vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/avatar.png"),
	}
	actor.Summary = vocab.NaturalLanguageValues{{Value: vocab.Content(game.Summary())}}
	actor.Attachment = vocab.ItemCollection{
		&PropertyValue{
			Type:  "PropertyValue",
			Name:  vocab.NaturalLanguageValues{{Value: vocab.Content("Created By")}},
			Value: vocab.NaturalLanguageValues{{Value: vocab.Content("<a href=\"https://chaos.social/@h4kor\" target=\"_blank\">@h4kor@chaos.social</a>")}},
		},
	}

	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
		jsonld.IRI(vocab.SecurityContextURI),
		jsonld.Context{
			jsonld.ContextElement{
				Term: "schema",
				IRI:  jsonld.IRI("http://schema.org#"),
			},
			jsonld.ContextElement{
				Term: "PropertyValue",
				IRI:  jsonld.IRI("schema:PropertyValue"),
			},
		},
	).Marshal(actor)
	if err != nil {
		slog.Error("Error marshalling", "err", err)
	}
	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}
