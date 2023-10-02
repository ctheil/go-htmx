package main

import (

	// locals

	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ctheil/go-htmx/postgres"
	"github.com/ctheil/go-htmx/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)

	if err := http.ListenAndServe(":42069", h); err != nil {
		log.Fatal(err)
	}

}
