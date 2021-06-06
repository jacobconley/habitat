package main

import (
	"net/http"

	"github.com/jacobconley/habitat/server"
	"github.com/rs/zerolog/log"
)

func main() { 
	log.Info().Msg("Starting test fixture server")

	router := server.Router()
	http.ListenAndServe(":3000", router)
}