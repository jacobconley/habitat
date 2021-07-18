package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jacobconley/habitat/server"
	"github.com/rs/zerolog/log"
)

func main() { 
	log.Info().Msg("Starting test fixture server")

	router := server.NewServer()

	router.Match("/test-get").
		GET().Raw(func(hab * server.Context) (result string, err error) {
			return "succ", nil
		})

	
	router.Match("/test-err").
		GET().Raw(func(hab * server.Context) (result string, err error) { 
			return "asdfl", fmt.Errorf("test error")
		})

	router.Match("/test-sql-404").
		GET().Raw(func(hab * server.Context) (result string, err error) {
			return "asdf", sql.ErrNoRows	
		})

	
	http.ListenAndServe(":3000", router)
}