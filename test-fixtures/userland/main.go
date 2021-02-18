package main

import (
	"habitat/src/server"

	log "github.com/sirupsen/logrus"
)

func main() { 
	log.Info("Starting test fixture server")

	if err := server.Serve(); err != nil { 
		log.Fatal("An error!")
		log.Fatal(err)
	}
}