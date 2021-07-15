package main

import (
	"net/http"

	"github.com/jacobconley/habitat/server"
	"github.com/rs/zerolog/log"
)

func main() { 
	log.Info().Msg("Starting test fixture server")

	router := server.NewServer()

	router.Match("/test-get").
		GET().Raw(func(hab server.Context) (result string, err error) {
			return "succ", nil
		})

	
	http.ListenAndServe(":3000", router)
}