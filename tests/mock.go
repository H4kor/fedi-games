package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/H4kor/fedi-games/web"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

type MockApServer struct {
	Server     *http.Server
	PrivateKey *rsa.PrivateKey
}

// / Mock server implementation for testing activity pub
func NewMockAPServer() MockApServer {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := privKey.Public().(*rsa.PublicKey)

	pubKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubKey),
		},
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		resource := r.URL.Query().Get("resource")
		parts := strings.Split(resource[5:], "@")
		req_name := parts[0]
		webfinger := web.WebfingerResponse{
			Subject: resource,

			Links: []web.WebfingerLink{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "http://localhost:7777/actors/" + req_name,
				},
			},
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		s, _ := json.Marshal(webfinger)
		w.Write(s)
	})
	mux.HandleFunc("GET /actors/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		actor := vocab.ServiceNew(vocab.IRI("http://localhost:7777/actors/" + name))
		actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(name)}}
		actor.Inbox = vocab.IRI("http://localhost:7777/actors/" + name + "/inbox")
		actor.PublicKey = vocab.PublicKey{
			ID:           vocab.ID("http://localhost:7777/actors/" + name + "#main-key"),
			Owner:        vocab.IRI("http://localhost:7777/actors/" + name),
			PublicKeyPem: string(pubKeyPem),
		}
		actor.Name = vocab.NaturalLanguageValues{{Value: vocab.Content(name)}}

		data, err := jsonld.WithContext(
			jsonld.IRI(vocab.ActivityBaseURI),
			jsonld.IRI(vocab.SecurityContextURI),
		).Marshal(actor)
		if err != nil {
			slog.Error("Error marshalling", "err", err)
		}
		w.Header().Add("Content-Type", "application/activity+json")
		w.Write(data)
	})

	srv := &http.Server{Addr: ":7777", Handler: mux}
	go func() {
		// always returns error. ErrServerClosed on graceful close
		println("Starting server on port 7777")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return MockApServer{
		Server:     srv,
		PrivateKey: privKey,
	}

}
