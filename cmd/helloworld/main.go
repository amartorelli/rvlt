package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/amartorelli/rvlt/pkg/api"
	"github.com/amartorelli/rvlt/pkg/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	// signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// config flags
	var lAddr = flag.String("listen", ":8080", "the address to listen on")
	var dbType = flag.String("db", "memory", "database type (memory/postgres)")
	var loglevel = flag.String("loglevel", "info", "log level (debug/info/warn/fatal)")
	flag.Parse()

	// log level
	switch *loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	db, err := database.NewDatabase(*dbType)
	if err != nil {
		log.Fatal(err)
	}

	a, err := api.NewHelloWorldAPI(*lAddr, db)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: start in a separate go routine
	go func() {
		err := a.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	s := <-stop
	log.Infof("received signal %s, gracefully shutting down", s.String())
	err = a.Stop()
	if err != nil {
		log.Error(err)
	}
}
