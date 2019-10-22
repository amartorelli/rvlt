package main

import (
	"github.com/amartorelli/rvlt/pkg/api"
	"github.com/amartorelli/rvlt/pkg/database"
	log "github.com/sirupsen/logrus"
)

func main() {

	db, err := database.NewDatabase("memory")
	if err != nil {
		log.Fatal(err)
	}

	a, err := api.NewHelloWorldAPI(":8080", db)
	if err != nil {
		log.Fatal(err)
	}
	a.Start()
}
