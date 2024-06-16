package main

import (
	"log"

	"github.com/H4kor/fedi-games/web"
)

func main() {
	server := web.DefaultServer()
	srv := server.Start()
	err := srv.ListenAndServe()

	log.Fatalf("%v", err)
}
