package main

import (
	"github.com/amartorelli/rvlt/pkg/api"
	"github.com/amartorelli/rvlt/pkg/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	// TODO:
	// read settings from env
	// handle signals

	db, err := database.NewDatabase("memory")
	if err != nil {
		log.Fatal(err)
	}

	a, err := api.NewHelloWorldAPI(":8080", db)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: start in a separate go routine
	a.Start()
}
