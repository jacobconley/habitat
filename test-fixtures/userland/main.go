package main

import (
	"net/http"

	"github.com/jacobconley/habitat/server"
	log "github.com/sirupsen/logrus"
)

func main() { 
	log.Info("Starting test fixture server")

	router := server.Router()
	http.ListenAndServe(":3000", router)
}