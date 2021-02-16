package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*
HabitatRouter returns the primary router for serving Habitat
	TODO: Include plugins
*/
func HabitatRouter() *mux.Router { 
	r := mux.NewRouter()

	// Maybe we should do query string instead?  seems cooler 
	r.HandleFunc("/!/assets/{path:[^:]+}:{hash:[a-fA-F0-9]+}", HandleAsset)

	return r
}

// Serve starts the server by invoking ListenAndServe
func Serve() error {  
	r := HabitatRouter()
	return http.ListenAndServe(":3000", r)
}
 
// // NewTestServer returns an httptest Server object with the HabitatRouter
// func NewTestServer() *httptest.Server { 
// 	log.Debug("Starting test server...") 
// 	r := HabitatRouter()
// 	return httptest.NewServer(r)
// }